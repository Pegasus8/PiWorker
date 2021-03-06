#!/usr/bin/env bash
#
# Installation script of PiWorker.
# https://github.com/Pegasus8/piworker

#
# ─── CONFIGS ────────────────────────────────────────────────────────────────────
#

CONFIG_SSL=1
CONFIG_INSTALL_DIR="$HOME/PiWorker"
CONFIG_INSTALL_SERVICE=1
CONFIG_AUTO_START_SERVICE=1
CONFIG_ADD_TO_PATH=1
CONFIG_AUTO_INSTALL_DEPENDENCIES=1
CONFIG_NEW_USER_STEP=1

#
# ─── INSTALL VARIABLES ──────────────────────────────────────────────────────────
#

LATEST_URL="https://api.github.com/repos/Pegasus8/piworker/releases/latest"

ARCH="" #$HOSTTYPE
OS="" #$OSTYPE
SUPPORTED_ARCH=("arm" "amd64")
SUPPORTED_OS=("linux")

PKG_MNGR=""
#APTGET_INSTALL_FLAGS=("install" "-y")
#PACMAN_INSTALL_FLAGS=("-Syu" "--noconfirm")

# Disable unicode.
LC_ALL=C
LANG=C

#
# ─── TEXT FORMATTING ─────────────────────────────────────────────────────────
#

BOLD='\e[1m'
#HIGHLIGHTED='\e[7m'
RED_FOREGROUND='\e[38;2;247;89;81m'
RESET='\e[m'

err() {
    echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')][ERROR]: $*" >&2
}

info() {
    echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')][INFO]: $*"
}

#
# ─── INSTALLATION FUNCTIONS ─────────────────────────────────────────────────────
#

GetOS() {
    PS3='> '
    echo "Select your OS:"
    select opt in "${SUPPORTED_OS[@]}"
    do
        case $opt in
            "${SUPPORTED_OS[0]}") OS="${SUPPORTED_OS[0]}"; break;;
            *) echo -e "${RED_FOREGROUND}Invalid option '$REPLY'${RESET}";;
        esac
    done
    printf "\n"
}

GetARCH() {
    PS3='> '
    echo "Select your ARCH:"
    select opt in "${SUPPORTED_ARCH[@]}"
    do
        case $opt in
            "${SUPPORTED_ARCH[0]}") ARCH="${SUPPORTED_ARCH[0]}"; break;;
            "${SUPPORTED_ARCH[1]}") ARCH="${SUPPORTED_ARCH[1]}"; break;;
            *) echo -e "${RED_FOREGROUND}Invalid option '$REPLY'${RESET}";;
        esac
    done
    printf "\n"
}

AddPath() {
    echo "export PATH=$PATH:$CONFIG_INSTALL_DIR/piworker" >>.bashrc
}

PrepareDirectory() {
    info "Checking if the directory of installation ('$CONFIG_INSTALL_DIR') already exists..."
    if [[ -d "$CONFIG_INSTALL_DIR" ]]; then
        err "The directory '$CONFIG_INSTALL_DIR' already exists"
        exit 1
    else
        if mkdir "$CONFIG_INSTALL_DIR"; then
            info "Directory of installation created correctly"
        else
            err "I can't create the directory where to install PiWorker. Maybe is an issue related with permissions?"
            exit 1
        fi
    fi
}

DownloadLatest() {
    workdir="$(mktemp -d)"
    trap 'rm -fr "$workdir"' RETURN
    local pwfile="pw.tar.gz"
    info "Starting the download..."

    local data
    local filename
    local url
    data="$(curl -sL "$LATEST_URL")"
    filename="$(echo "$data" | grep piworker-"$OS"_"$ARCH"- | grep name | head -1 | cut -d \" -f 4)"
    url="$(echo "$data" | grep piworker-"$OS"_"$ARCH"- | grep browser_download_url | head -1 | cut -d \" -f 4)"

    info "Downloading '$filename'..."
    if ! curl \
        --request GET \
        -sL \
        --url "$url" \
        --output "$workdir/$pwfile"
    then
        err "Download of '$filename' failed"
        exit 1
    fi
    info "'$filename' downloaded!"

    info "Downloading the sha256sum of '$filename'..."
    if ! curl \
        --request GET \
        -sL \
        --url "$url.sha256sum" \
        --output "$workdir/$pwfile.sha256sum"
    then
        err "Download of '$filename.sha256sum' failed"
        exit 1
    fi
    info "'$filename.sha256sum' downloaded!"

    local sha256sum_downloaded
    local sha256sum_calc
    sha256sum_downloaded="$(< "$workdir/$pwfile.sha256sum" awk '{print $1}')"
    sha256sum_calc="$(sha256sum "$workdir/$pwfile" | awk '{print $1}')"

    if [[ $sha256sum_downloaded == "$sha256sum_calc" ]]; then
        info "sha256sum of the file '$filename' verified correctly! ($sha256sum_downloaded)"
    else
        err "The downloaded sha256sum ($sha256sum_downloaded) of the file '$filename' does not match with the calculated one ($sha256sum_calc)"
        exit 1
    fi

    local file_dir="$workdir/PiWorker"

    if ! mkdir "$file_dir"; then
        err "Error when trying to make the directory ('$file_dir') where decompress the file '$pwfile'"
        exit 1
    fi

    if ! tar -C "$file_dir" -xzf "$workdir/$pwfile"; then
        err "Error when trying to decompress the file '$workdir/$pwfile' in '$file_dir'"
        exit 1
    fi

    if [[ -x "$file_dir/piworker" ]]; then
        if mv -f "$file_dir/piworker" "$CONFIG_INSTALL_DIR"; then
            info "Executable moved from '$file_dir/piworker' to '$CONFIG_INSTALL_DIR/piworker'"
        else
            err "Error when trying to move the executable from '$file_dir/piworker' to '$CONFIG_INSTALL_DIR/piworker'"
            exit 1
        fi
    else
        err "I can't find the executable in '$file_dir'"
        exit 1
    fi
}

InstallService() {
    info "Installing the service (with superuser permissions)..."
    
    if sudo -u root "$CONFIG_INSTALL_DIR/piworker" --service install; then 
        info "Service installed!"
        if [[ "$CONFIG_AUTO_START_SERVICE" -eq 0 ]]; then
            return
        fi

        info "Starting the service..."
        if sudo -u root "$CONFIG_INSTALL_DIR/piworker" --service start; then
            info "Service started!"
        else
            err "Cannot start the service. You can do it manually executing: $CONFIG_INSTALL_DIR/piworker --service start"
        fi
    else 
        err "The service was not installed. You can try to do it manually running: $CONFIG_INSTALL_DIR/piworker --service install"
    fi

    info "Making sure that the user '$USER' has the ownership of the directory '$CONFIG_INSTALL_DIR'"
    if sudo -u root chown -R "$USER" "$CONFIG_INSTALL_DIR"; then
        info "Done!"
    else
        err "Can't change the ownership of the directory '$CONFIG_INSTALL_DIR'. If you want to do it manually run: sudo -u root chown -R $USER $CONFIG_INSTALL_DIR"
    fi

    info "Making sure that the user '$USER' has the permissions of the directory '$CONFIG_INSTALL_DIR'"
    if sudo -u root chmod -R u+rwx "$CONFIG_INSTALL_DIR"; then
        info "Done!"
    else
        err "Can't change the permissions of the directory '$CONFIG_INSTALL_DIR'. If you want to do it manually run: sudo -u root chmod -R u+rwx $CONFIG_INSTALL_DIR"
    fi
}

GenerateOpenSSLCertificate() {
    info "Generating a new self signed certificate..."
    if openssl req \
        -subj '/O=PiWorker' \
        -new \
        -newkey \
        rsa:2048 \
        -sha256 \
        -days 365 \
        -nodes \
        -x509 \
        -keyout "$CONFIG_INSTALL_DIR/server.key" \
        -out "$CONFIG_INSTALL_DIR/server.crt"
    then
        info "SSL certificates generated successfully"
        return 0
    else
        err "The generation of SSL certificates has failed with status $?"
        return 1
    fi
}

DefinePackageManager() {
    info "Trying to differentiate the package manager of the system..."
    if hash apt-get 2>/dev/null; then
        PKG_MNGR="apt-get"
        info "$PKG_MNGR identified"
    elif hash pacman 2>/dev/null; then
        PKG_MNGR="pacman"
        info "$PKG_MNGR identified"
    else
        err "Cannot identify the package manager of the system"
        read -r -p "Should I continue without checking the dependencies? (Y/N): " confirm && [[ $confirm == [yY] || $confirm == [yY][eE][sS] ]] || exit 1
    fi
}

InstallDependencies() {
    info "Checking dependencies..."
    deps=("curl" "grep")
    if [[ $CONFIG_SSL -ne 0 ]]; then
        deps+=("openssl")
    fi
    local aptget_pkgmngr="apt-get"
    local pacman_pkgmngr="pacman"
    local aptget_updated=0

    for dependency in "${deps[@]}"; do
        if hash "$dependency" 2>/dev/null; then
            info "$dependency already installed!"
        else
            info "$dependency not installed, installing it (with superuser permissions)..."
            if [[ $PKG_MNGR -eq "$aptget_pkgmngr" ]]; then
                if [[ $aptget_updated -eq 0 ]]; then
                    info "Updating apt-get cache for first time (with superuser permissions)..."
                    if sudo -u root apt-get update; then
                        aptget_updated=1
                        info "apt-get cache updated!"
                    else
                        err "Cannot update the cache of apt-get"
                    fi
                fi

                if ! sudo -u root apt-get install -y "$dependency"; then
                    err "Error when trying to install $dependency"
                    return 1
                fi

                info "$dependency installed!"
            elif [[ $PKG_MNGR -eq "$pacman_pkgmngr" ]]; then
                if ! sudo -u root pacman -Sy --noconfirm "$dependency"; then
                    err "Error when trying to install $dependency"
                    return 1
                fi

                info "$dependency installed!"
            else
                return 1
            fi
        fi
    done

    return 0
}

NewUser() {
    info "Starting the creation of the first user..."
    if ! cd "$CONFIG_INSTALL_DIR"; then
        err "I cannot enter to the directory '$CONFIG_INSTALL_DIR' to create a new user"
        return
    fi

    while true; do
        read -r -p "Username: " username
        read -r -p "Password: " password
        read -r -p "Administrator? (Y/N): " admin
        local admin_str=""
        if [[ $admin == [yY] || $admin == [yY][eE][sS] ]]; then
            admin=1
            admin_str="yes"
        else
            admin=0
            admin_str="no"
        fi
        printf "\n"

        echo -e "Let's review the data before creating the user"
        echo -e "Username choosed: ${BOLD}$username${RESET}"
        echo -e "Password choosed: ${BOLD}$password${RESET}"
        echo -e "Admin: ${BOLD}$admin_str${RESET}"
        read -r -p "Is that info correct? (Y/N): " confirm
        if [[ $confirm == [yY] || $confirm == [yY][eE][sS] ]]; then
            if [[ $admin -ne 0 ]]; then
                ./piworker --new-user --username "$username" --password "$password" --admin
            else
                ./piworker --new-user --username "$username" --password "$password"
            fi
            break
        else
            echo "No problem, let's fill the information again"
        fi
    done
}

#
# ─── EXECUTION ──────────────────────────────────────────────────────────────────
#

echo "PiWorker installer"
echo "------------------"
echo -e "GitHub repo: https://github.com/Pegasus8/piworker\n"

# Identify the characteristics of the host to choose the apropiate binary.
GetOS
GetARCH
echo "---------------------"
echo -e "Selected OS: ${BOLD}$OS${RESET}"
echo -e "Selected ARCH: ${BOLD}$ARCH${RESET}"
read -r -p "Is that info correct? (Y/N): " confirm && [[ $confirm == [yY] || $confirm == [yY][eE][sS] ]] || exit 1
echo "Nice, let's continue!"

# if [[ ! " ${SUPPORTED_ARCHS[@]} " =~ " ${ARCH} " ]]; then
#     print_redf "For now, PiWorker doesn't support your architecture ($ARCH). Sorry."
#     print_blueb "If you want PiWorker to support your architecture you can open an issue in the Github repository: https://github.com/Pegasus8/piworker/issues/new"
#     exit 1
# fi

PrepareDirectory

if [[ $CONFIG_AUTO_INSTALL_DEPENDENCIES -ne 0 ]]; then
    DefinePackageManager
    if InstallDependencies; then
        if [[ $CONFIG_SSL -ne 0 ]]; then
            if ! GenerateOpenSSLCertificate; then
                read -r -p "Cannot generate a self signed SSL certificate, should I continue? (Y/N): " confirm && [[ $confirm == [yY] || $confirm == [yY][eE][sS] ]] || exit 1
            fi
        fi
    else
        exit 1
    fi
fi

DownloadLatest

if [[ $CONFIG_ADD_TO_PATH -ne 0 ]]; then
    AddPath
fi

if [[ $CONFIG_NEW_USER_STEP -ne 0 ]]; then
    NewUser
fi

if [[ $CONFIG_INSTALL_SERVICE -ne 0 ]]; then
    InstallService
fi

info "Installation finished!"
