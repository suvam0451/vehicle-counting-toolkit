

## Generates batches of images from video files in "input" folder

$VideoFiles = Get-Childitem "input" -File | Where-Object {$_.extension -eq ".mp4"}

foreach ($VideoFile in $VideoFiles) {
    # Always confirm that the "intermediate" directory is present
    if(!(Test-Path -Path "intermediate")) {
        New-Item -ItemType Directory -Force -Path "intermediate"
    }
    if(!(Test-Path -Path "imageset")) {
        New-Item -ItemType Directory -Force -Path "imageset"
    }
    $VideoName = [io.path]::GetFileNameWithoutExtension($VideoFile)

    # generate folders to an intermediate directory
    $Name ="imagesets/" + [io.path]::GetFileNameWithoutExtension($VideoFile)
    # Create directory for each video file, if missing
    if(!(Test-Path -Path $Name)) {
        New-Item -ItemType Directory -Force -Path $Name
    }
    # Run yolo_mark for each video
    ./bin/yolo_mark.exe $Name cap_video $VideoFile 10

    # Generate the image list for folders
    $Images = Get-Childitem $Name

    # $Images |get-member
    $List = $Images | Where-Object {$_.extension -eq ".jpg"}
    $List |Format-Table fullname -HideTableHeaders |out-file $Name"\output.txt"
    # Get rid of above-most blankline
    Get-Content $Name"\output.txt" | Where-Object {$_ -ne ""} | out-file "intermediate\$VideoName.txt"
    Remove-Item $Name"\output.txt"
    # $List | format-table name
}