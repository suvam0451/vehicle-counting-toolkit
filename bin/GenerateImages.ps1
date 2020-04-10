$VideoFiles = Get-Childitem . -File | Where-Object {$_.extension -eq ".mp4"}

foreach ($VideoFile in $VideoFiles) {
    $Name =[io.path]::GetFileNameWithoutExtension($VideoFile)
    # Create directory for each video file, if missing
    if(!(Test-Path -Path $Name)) {
        New-Item -ItemType Directory -Force -Path $Name
    }
    # Run yolo_mark for each video
    ./yolo_mark.exe $Name cap_video $VideoFile 10

    # Generate the image list for folders
    $Images = Get-Childitem $Name

    # $Images |get-member
    $List = $Images | Where-Object {$_.extension -eq ".jpg"}
    $List |Format-Table fullname -HideTableHeaders |out-file $Name"\output.txt"
    # Get rid of above-most blankline
    Get-Content $Name"\output.txt" | Where-Object {$_ -ne ""} | out-file $Name"\list.txt"
    Remove-Item $Name"\output.txt"
    # $List | format-table name
}