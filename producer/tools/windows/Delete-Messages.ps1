$response = Invoke-WebRequest -Uri 'http://localhost:8080/clearMessages' -Method 'GET'
Write-Output $response.Content
