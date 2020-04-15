# $Folders = Get-Childitem . -Directory | Select-Object FullName
$TextFiles = Get-Childitem "intermediate" -File | Where-Object {$_.extension -eq ".txt"} | Select-Object FullName

# For each supported text files...
foreach ($TextFile in $TextFiles) {
    # echo $TextFile/"list.txt"
    $FolderPath = $TextFile.FullName
    $Filename = [io.path]::GetFileNameWithoutExtension($FolderPath)

    # $InputPath = $FolderPath + "\list.txt"
    $OutputPath = "output/" + $Filename + ".json"
    
    # Run darknet on that file
    Get-Content $FolderPath | ./bin/darknet.exe detector test data/coco.data cfg/yolov3.cfg cfg/yolov3.weights -dont_show -ext_output -out $OutputPath
    # Rename-Item -Path "out.json" -NewName $Outputpath
    # Copy-Item $Outputpath -Destination "output"
}