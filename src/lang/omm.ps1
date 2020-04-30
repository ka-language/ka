#!/usr/bin/env pwsh
$basedir = Split-Path $MyInvocation.MyCommand.Definition -Parent
$calldir = $PSScriptRoot

#call the actual script
& "$basedir/lang.exe" "$calldir\" $args
exit 0