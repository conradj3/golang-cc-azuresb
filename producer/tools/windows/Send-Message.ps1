function Send-Message {
    param (
        [Parameter(Mandatory = $true)]
        [int]$Count
    )

    $body = @{ count = $Count } | ConvertTo-Json
    Invoke-WebRequest -Uri 'http://localhost:8080/createMessages' -Method 'POST' -Body $body -ContentType 'application/json'

}

# Example usage:
Send-Message -Count 300
