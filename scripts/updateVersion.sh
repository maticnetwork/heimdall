#!/bin/bash
set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

echo "The version is of form - VersionMajor.VersionMinor.VersionPatch-VersionMeta"
echo "Let's take 0.3.4-beta as an example. Here:"
echo "* VersionMajor is - 0"
echo "* VersionMinor is - 3"
echo "* VersionPatch is - 4"
echo "* VersionMeta is  - beta"
echo ""
echo "Now, enter the new version step-by-step below:"

version=""

# VersionMajor
read -p "* VersionMajor: " VersionMajor
if [ -z "$VersionMajor" ]
then
    echo "VersionMajor cannot be NULL"
    exit -1
fi
version+=$VersionMajor

# VersionMinor
read -p "* VersionMinor: " VersionMinor
if [ -z "$VersionMinor" ]
then
    echo "VersionMinor cannot be NULL"
    exit -1
fi
version+="."$VersionMinor

# VersionPatch
read -p "* VersionPatch: " VersionPatch
if [ -z "$VersionPatch" ]
then
    echo "VersionPatch cannot be NULL"
    exit -1
fi
version+="."$VersionPatch

# VersionMeta (optional)
read -p "* VersionMeta (optional, press enter if not needed): " VersionMeta
if [[ ! -z "$VersionMeta" ]]
then
    version+="-"$VersionMeta
fi

echo ""
echo "New version is: $version"

# update version in all the 6 templates and 1 deb control file
replaceVersion="Version: "$version
replaceStandards="Standards-Version: v"$version
fileArray=(
    "${DIR}/../packaging/deb/heimdalld/DEBIAN/control"
    "${DIR}/../packaging/templates/package_scripts/control"
    "${DIR}/../packaging/templates/package_scripts/control.arm64"
    "${DIR}/../packaging/templates/package_scripts/control.profile.amd64"
    "${DIR}/../packaging/templates/package_scripts/control.profile.arm64"
    "${DIR}/../packaging/templates/package_scripts/control.validator"
    "${DIR}/../packaging/templates/package_scripts/control.validator.arm64"
)
for file in ${fileArray[@]}; do
    # get the line starting with `Version` in the control file and store it in the $tempVersion variable
    tempVersion=$(grep "^Version.*" $file)
    sed -i '' "s%$tempVersion%$replaceVersion%" $file
done

fileArrayStandards=(
    "${DIR}/../packaging/deb/heimdalld/DEBIAN/control"
    "${DIR}/../packaging/templates/package_scripts/control.validator"
    "${DIR}/../packaging/templates/package_scripts/control.validator.arm64"
)
for file in ${fileArrayStandards[@]}; do
    # get the line starting with `Standards-Version` in the control file and store it in the $tempStandards variable
    tempStandards=$(grep "^Standards-Version.*" $file)
    sed -i '' "s%$tempStandards%$replaceStandards%" $file
done

echo ""
echo "Updating Version Done"

exit 0
