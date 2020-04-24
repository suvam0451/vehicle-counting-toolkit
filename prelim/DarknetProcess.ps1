$Folders = Get-Childitem "intermediate" -File | Where-Object { $_.extension -eq ".txt" }

foreach ($ListFile in $Folders) {
    $ListFileName = [io.path]::GetFileNameWithoutExtension($ListFile.FullName)
    $OutputPath = "output/$ListFileName.json"
    Get-Content $ListFile | .\darknet.exe detector test cfg/coco.data cfg/yolov3.cfg cfg/yolov3.weights -dont_show -out $OutputPath
    # Replace "\" with "/" to avoid JSON parsing errors across platforms
    ((Get-Content -Path $OutputPath -Raw) -replace '\\', '/') | Set-Content $OutputPath

    # This is the root array for the frames
    $RootJSONObject = @();

    $JSONFromFile = Get-Content -Raw -Path $OutputPath | ConvertFrom-JSON
    foreach ($JSONFrame in $JSONFromFile) {
        # This is the root for each frame object
        $base = @{frame_id = $JSONFrame.frame_id; objects = @() }

        foreach ($ObjectEntry in $JSONFrame.objects) {
            $ID = $ObjectEntry.class_id
            $XCenter = $ObjectEntry.relative_coordinates.center_x
            $YCenter = $ObjectEntry.relative_coordinates.center_y
            $Confidence = $ObjectEntry.confidence
            if ($ID -eq 0) {
                # Skipping "person"
            } else {
                $base.objects += [pscustomobject]@{class_id=$ID;center_x=$XCenter;center_y=$YCenter;confidence=$Confidence;}
            }
        }
        $RootJSONObject += $base
    }
    $RootJSONObject | ConvertTo-Json -depth 100 | Out-File "LonelyTogether.json"
}