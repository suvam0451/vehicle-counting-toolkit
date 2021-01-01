$Dirs = Get-Childitem . -Recurse -Directory
foreach ($Dir in $Dirs) {
    $name = Join-Path -Path $Dir.FullName -ChildPath ""

    # List all images in the folder
    $Subdir = Get-Childitem $name -Recurse -File
    $List = $Subdir | where {$_.extension -eq ".jpg"}

    $OutputFile = $Dir.Name + ".txt"
    $List |ft FullName -HideTableHeaders |out-file $OutputFile
    # $List | format-table name
}
# Copy-Item "yolo_mark.exe" -Destination "copyto"