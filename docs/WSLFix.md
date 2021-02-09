# Fix WSL 1 to Update to WSL 2

These instructions assume that Windows 10 is fully updated on the system.

Before starting, complete the [Docker instructions.](DockerWindows.md)
There is a chance that an error might occur while running the "wsl --set-default-version 2" command.


1. Follow the website listed in the error message. On the website, there will be a .msi update package to download.
1. Following this download, run this command in PowerShell:

```powershell
wsl --set-default-version 2
```

There is a chance that after running this, you will receive the same error from before. If that occurs follow the steps below.

1. Open PowerShell as an Administrator. (This can be completed by right clicking on the PowerShell icon and selecting: "Run as Administrator")

2. Run this command:

```powershell
dism.exe /online /enable-feature /featurename:Microsoft-Windows-Subsystem-Linux /all /norestart
```

3. After that command is completed, run:

```powershell
dism.exe /online /enable-feature /featurename:VirtualMachinePlatform /all /norestart
```

4. Retry the "set default version" command from the beginning steps..
5. Restart your system.
6. The system should be updated and you should now be able to continue with the Docker instructions.
