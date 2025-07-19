#!/usr/bin/env pwsh
param(
  [String]$Version = "latest",
  # Skips adding the autoprovisioner.exe directory to the user's %PATH%
  [Switch]$NoPathUpdate = $false,
  # Skips adding autoprovisioner to the list of installed programs
  [Switch]$NoRegisterInstallation = $false,

  # Debugging: Always download with 'Invoke-RestMethod' instead of 'curl.exe'
  [Switch]$DownloadWithoutCurl = $false
);

# filter out 32 bit + ARM
if (-not ((Get-CimInstance Win32_ComputerSystem)).SystemType -match "x64-based") {
  Write-Output "Install Failed:"
  Write-Output "Autoprovisioner for Windows is currently only available for x86 64-bit Windows.`n"
  return 1
}

$ErrorActionPreference = "Stop"

# These three environment functions are roughly copied from https://github.com/prefix-dev/pixi/pull/692
# They are used instead of `SetEnvironmentVariable` because of unwanted variable expansions.
function Publish-Env {
  if (-not ("Win32.NativeMethods" -as [Type])) {
    Add-Type -Namespace Win32 -Name NativeMethods -MemberDefinition @"
[DllImport("user32.dll", SetLastError = true, CharSet = CharSet.Auto)]
public static extern IntPtr SendMessageTimeout(
    IntPtr hWnd, uint Msg, UIntPtr wParam, string lParam,
    uint fuFlags, uint uTimeout, out UIntPtr lpdwResult);
"@
  }
  $HWND_BROADCAST = [IntPtr] 0xffff
  $WM_SETTINGCHANGE = 0x1a
  $result = [UIntPtr]::Zero
  [Win32.NativeMethods]::SendMessageTimeout($HWND_BROADCAST,
    $WM_SETTINGCHANGE,
    [UIntPtr]::Zero,
    "Environment",
    2,
    5000,
    [ref] $result
  ) | Out-Null
}

function Write-Env {
  param([String]$Key, [String]$Value)

  $RegisterKey = Get-Item -Path 'HKCU:'

  $EnvRegisterKey = $RegisterKey.OpenSubKey('Environment', $true)
  if ($null -eq $Value) {
    $EnvRegisterKey.DeleteValue($Key)
  } else {
    $RegistryValueKind = if ($Value.Contains('%')) {
      [Microsoft.Win32.RegistryValueKind]::ExpandString
    } elseif ($EnvRegisterKey.GetValue($Key)) {
      $EnvRegisterKey.GetValueKind($Key)
    } else {
      [Microsoft.Win32.RegistryValueKind]::String
    }
    $EnvRegisterKey.SetValue($Key, $Value, $RegistryValueKind)
  }

  Publish-Env
}

function Get-Env {
  param([String] $Key)

  $RegisterKey = Get-Item -Path 'HKCU:'
  $EnvRegisterKey = $RegisterKey.OpenSubKey('Environment')
  $EnvRegisterKey.GetValue($Key, $null, [Microsoft.Win32.RegistryValueOptions]::DoNotExpandEnvironmentNames)
}

# The installation of autoprovisioner is it's own function for better error handling
function Install-Autoprovisioner {
  param(
    [string]$Version
  );

  # if a semver is given, we need to adjust it to this format: v0.0.0
  if ($Version -match "^\d+\.\d+\.\d+$") {
    $Version = "v$Version"
  }

  $Arch = "x64"

  $AutoprovisionerRoot = if ($env:AUTOPROVISIONER_INSTALL) { $env:AUTOPROVISIONER_INSTALL } else { "${Home}\.autoprovisioner" }
  $AutoprovisionerBin = mkdir -Force "${AutoprovisionerRoot}\bin"

  try {
    Remove-Item "${AutoprovisionerBin}\autoprovisioner.exe" -Force
  } catch [System.Management.Automation.ItemNotFoundException] {
    # ignore
  } catch [System.UnauthorizedAccessException] {
    $openProcesses = Get-Process -Name autoprovisioner | Where-Object { $_.Path -eq "${AutoprovisionerBin}\autoprovisioner.exe" }
    if ($openProcesses.Count -gt 0) {
      Write-Output "Install Failed - An older installation exists and is open. Please close open Autoprovisioner processes and try again."
      return 1
    }
    Write-Output "Install Failed - An unknown error occurred while trying to remove the existing installation"
    Write-Output $_
    return 1
  } catch {
    Write-Output "Install Failed - An unknown error occurred while trying to remove the existing installation"
    Write-Output $_
    return 1
  }

  $Target = "autoprovisioner-windows-$Arch"
  $BaseURL = "https://github.com/zerosync-co/autoprovisioner-cli/releases"
  $URL = "$BaseURL/$(if ($Version -eq "latest") { "latest/download" } else { "download/$Version" })/$Target.zip"

  $ZipPath = "${AutoprovisionerBin}\$Target.zip"

  $DisplayVersion = $(
    if ($Version -eq "latest") { "Autoprovisioner" }
    elseif ($Version -match "^v\d+\.\d+\.\d+$") { "Autoprovisioner $($Version.Substring(1))" }
    else { "Autoprovisioner tag='${Version}'" }
  )

  $null = mkdir -Force $AutoprovisionerBin
  Remove-Item -Force $ZipPath -ErrorAction SilentlyContinue

  # curl.exe is faster than PowerShell 5's 'Invoke-WebRequest'
  # note: 'curl' is an alias to 'Invoke-WebRequest'. so the exe suffix is required
  if (-not $DownloadWithoutCurl) {
    curl.exe "-#SfLo" "$ZipPath" "$URL"
  }
  if ($DownloadWithoutCurl -or ($LASTEXITCODE -ne 0)) {
    Write-Warning "The command 'curl.exe $URL -o $ZipPath' exited with code ${LASTEXITCODE}`nTrying an alternative download method..."
    try {
      # Use Invoke-RestMethod instead of Invoke-WebRequest because Invoke-WebRequest breaks on
      # some machines
      Invoke-RestMethod -Uri $URL -OutFile $ZipPath
    } catch {
      Write-Output "Install Failed - could not download $URL"
      Write-Output "The command 'Invoke-RestMethod $URL -OutFile $ZipPath' failed`n"
      return 1
    }
  }

  if (!(Test-Path $ZipPath)) {
    Write-Output "Install Failed - could not download $URL"
    Write-Output "The file '$ZipPath' does not exist. Did an antivirus delete it?`n"
    return 1
  }

  try {
    $lastProgressPreference = $global:ProgressPreference
    $global:ProgressPreference = 'SilentlyContinue';
    Expand-Archive "$ZipPath" "$AutoprovisionerBin" -Force
    $global:ProgressPreference = $lastProgressPreference
    if (!(Test-Path "${AutoprovisionerBin}\autoprovisioner.exe")) {
      throw "The file '${AutoprovisionerBin}\autoprovisioner.exe' does not exist. Download is corrupt or intercepted by Antivirus?`n"
    }
  } catch {
    Write-Output "Install Failed - could not unzip $ZipPath"
    Write-Error $_
    return 1
  }

  Remove-Item $ZipPath -Force

  # Test that the binary works
  $AutoprovisionerVersion = "$(& "${AutoprovisionerBin}\autoprovisioner.exe" --version 2>&1)"
  if ($LASTEXITCODE -ne 0) {
    Write-Output "Install Failed - could not verify autoprovisioner.exe"
    Write-Output "The command '${AutoprovisionerBin}\autoprovisioner.exe --version' exited with code ${LASTEXITCODE}"
    Write-Output "Output: $AutoprovisionerVersion`n"
    return 1
  }

  $C_RESET = [char]27 + "[0m"
  $C_GREEN = [char]27 + "[1;32m"

  Write-Output "${C_GREEN}Autoprovisioner was installed successfully!${C_RESET}"
  Write-Output "The binary is located at ${AutoprovisionerBin}\autoprovisioner.exe`n"

  $hasExistingOther = $false;
  try {
    $existing = Get-Command autoprovisioner -ErrorAction Stop
    if ($existing.Source -ne "${AutoprovisionerBin}\autoprovisioner.exe") {
      Write-Warning "Note: Another autoprovisioner.exe is already in %PATH% at $($existing.Source)`nTyping 'autoprovisioner' in your terminal will not use what was just installed.`n"
      $hasExistingOther = $true;
    }
  } catch {}

  if (-not $NoRegisterInstallation) {
    $rootKey = $null
    try {
      $RegistryKey = "HKCU:\Software\Microsoft\Windows\CurrentVersion\Uninstall\Autoprovisioner"
      $rootKey = New-Item -Path $RegistryKey -Force
      New-ItemProperty -Path $RegistryKey -Name "DisplayName" -Value "Autoprovisioner" -PropertyType String -Force | Out-Null
      New-ItemProperty -Path $RegistryKey -Name "InstallLocation" -Value "${AutoprovisionerRoot}" -PropertyType String -Force | Out-Null
      New-ItemProperty -Path $RegistryKey -Name "DisplayIcon" -Value $AutoprovisionerBin\autoprovisioner.exe -PropertyType String -Force | Out-Null
      New-ItemProperty -Path $RegistryKey -Name "UninstallString" -Value "powershell -c `"Remove-Item -Recurse -Force '$AutoprovisionerRoot'; [Environment]::SetEnvironmentVariable('Path', ([Environment]::GetEnvironmentVariable('Path', 'User') -split ';' | Where-Object { `$_ -ne '$AutoprovisionerBin' }) -join ';', 'User')`"" -PropertyType String -Force | Out-Null
    } catch {
      if ($rootKey -ne $null) {
        Remove-Item -Path $RegistryKey -Force
      }
    }
  }

  if(!$hasExistingOther) {
    # Only try adding to path if there isn't already an autoprovisioner.exe in the path
    $Path = (Get-Env -Key "Path") -split ';'
    if ($Path -notcontains $AutoprovisionerBin) {
      if (-not $NoPathUpdate) {
        $Path += $AutoprovisionerBin
        Write-Env -Key 'Path' -Value ($Path -join ';')
        $env:PATH = $Path;
      } else {
        Write-Output "Skipping adding '${AutoprovisionerBin}' to the user's %PATH%`n"
      }
    }

    Write-Output "To get started, restart your terminal/editor, then type `"autoprovisioner`"`n"
  }

  $LASTEXITCODE = 0;
}

Install-Autoprovisioner -Version $Version
