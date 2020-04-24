## INFO : Generates batches of images from video files in "input" folder

# Parameter inputs
$FrameSkip = 100

# We should run the scripts from inside
Set-Location bin
$VideoFiles = Get-Childitem "../input" -File | Where-Object {$_.extension -eq ".mp4" -or $_.extension -eq ".mpg" -or $_.extension -eq ".MTS" -or $_.extension -eq ".avi"}

foreach ($VideoFile in $VideoFiles) {
    # Always confirm that the "intermediate" directory is present
    if(!(Test-Path -Path "../intermediate")) {
        New-Item -ItemType Directory -Force -Path "intermediate"
    }
    if(!(Test-Path -Path "../imagesets")) {
        New-Item -ItemType Directory -Force -Path "imagesets"
    }
    $VideoName = [io.path]::GetFileNameWithoutExtension($VideoFile)
    # generate folders to an intermediate directory
    $Name ="../imagesets/" + [io.path]::GetFileNameWithoutExtension($VideoFile)
    # Create directory for each video file, if missing
    if(!(Test-Path -Path $Name)) {
        New-Item -ItemType Directory -Force -Path $Name
    }
    # Run yolo_mark for each video
    ./yolo_mark.exe $Name cap_video $VideoFile $FrameSkip

    # Generate the image list for folders
    $Images = Get-Childitem $Name

    # $Images |get-member
    $List = $Images | Where-Object {$_.extension -eq ".jpg"}
    $List |Format-Table fullname -HideTableHeaders |out-file $Name"\output.txt"
    # Get rid of above-most blankline
    Get-Content $Name"\output.txt" | Where-Object {$_ -ne ""} | out-file "../intermediate\$VideoName.txt"
    Remove-Item $Name"\output.txt"
}

# Exit directory and notify completion
Set-Location ..
Write-Output "Done..."
