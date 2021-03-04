# Launching the cluster on Windows

These instructions assume WSL 2 and Docker Desktop are already installed.

1. Delete the repository from your local computer.
1. Open Docker Desktop and **DELETE ALL CONTAINERS!**
1. Install WSL 2 by installing Ubuntu from Microsoft Store if you have not
   already.
1. Open a Ubuntu WSL terminal.
1. You may be asked to make a username and password if you have not already.
   Make a note of these.
1. Run the command `cd` to get back to your WSL home directory.
1. If the path starts with `/mnt/c` and not `~`, stop here and ask for help
   because something isn't right.
1. Clone the repository:

   ```sh
    (
        export SERVER="dev.azure.com";
        export ORG="S21-Team14-Godfrey-Brazil-Caples-Sharpe";
        export REPO="S21-Team14-Godfrey.Brazil.Caples.Sharpe";

        git clone https://${ORG}@${SERVER}/${ORG}/${REPO}/_git/${REPO} cpsc491;
    )
   ```

1. Open Docker Desktop.
1. Under Settings > Resources > WSL Integration, turn the switch ON next to
   Ubuntu.
1. At the Ubuntu terminal, run this command:

   ```sh
   sudo apt update && sudo apt install make
   ```

1. **REBOOT THE ENTIRE COMPUTER.** (seriously, don't skip this!)
1. Wait for Docker Desktop to finish launching.
1. Open a Ubuntu WSL terminal.
1. Enter command `cd ~/cpsc491`.
1. Enter command `make`.
1. Wait for containers to start.
1. Open VSCode.
1. Install extension "Remote - WSL".
1. Click green >< button in lower left of window.
1. Select "Remote WSL - Open folder in WSL".
1. You will likely start in the `docker-desktop` folder. Go up one directory.
1. Enter the `ubuntu` directory.
1. Select the folder at the path `home/YOUR_USERNAME/cpsc491`.
1. Wait for the folder to open (will take several minutes the first time).
1. Visit [http://localhost:8000](http://localhost:8000) in a browser.
1. Find `web/src/App.js` and make an edit and SAVE.
1. Wait for the app to recompile. You can watch the terminal window to see
   progress.
1. Watch the web browser and profit!

You may need to reinstall VSCode extensions inside WSL after these steps are
complete.
