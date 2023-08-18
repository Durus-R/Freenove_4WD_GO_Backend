Param(
    [Parameter(Position=0)]
    $param1
)

$ErrorActionPreference = "Stop"

function protoc  {
    New-Item -ItemType Directory -Force -Path dist

    & protoc --go_out=dist/ `
  --go_opt=paths=source_relative `
  --go-grpc_out=dist/ `
  --go-grpc_opt=paths=source_relative proto/*.proto
}

$TARGET_ARCH_CROSS = "linux/arm64"
$BUILD_COMMAND = "go build"
$BUILD_TAG = "-tags pi"
$APP_NAME = "Freenove_4WD_GO_Backend"
$DOCKER_EXEC_COMMAND = "docker run --rm -v `"$PWD`":/usr/src/$APP_NAME --platform $TARGET_ARCH_CROSS
    -w /usr/src/$APP_NAME go-cross-builder:latest"


function test_docker {
    try {
        $info = docker info
        Write-Host "Info of docker:"
        Write-Host $info
        return 0
    }
    catch {
        Write-Host "Docker is not running. Please start the Docker engine or Docker Desktop."
        return 1
    }
}

function test_go {
    # Get go version

    try {
        $goCommand = Get-Command go -ErrorAction Stop
        Write-Output "Found Go: $($goCommand.Source)"
    } catch {
        Write-Output "Go is not on the system PATH"
        return 1
    }

    $goVersion = go versiontry {
        $goCommand = Get-Command go -ErrorAction Stop
        Write-Output "Go is installed in $($goCommand.Source)"
    } catch {
        Write-Output "Go is not on the system PATH"
    }
    $goVersion = $goVersion.Split(' ')[2].Replace("go", "")
    $goVersionMain = $goVersion.Split('.')[0..1] -join '.'

    # Specify the version to check
    $desiredVersion = "1.20"

    # Check if the installed go version is greater than or equal to the desired version
    if ([version]$goVersionMain -ge [version]$desiredVersion) {
        Write-Output "Go version is compatible with $desiredVersion. The installed version is: $goVersionMain"
        return 0
    } else {
        Write-Output "Go version is incompatible with $desiredVersion. The installed version is: $goVersionMain"
        return 1
    }
}

if($param1 -eq 'cross') {
    $ret = test_docker
    if ($ret -eq 1) {
        exit 1
    }
    & docker buildx build --platform $TARGET_ARCH_CROSS --tag go-cross-builder .
    & $DOCKER_EXEC_COMMAND bash protoc.sh
    & $DOCKER_EXEC_COMMAND go build $BUILD_TAG -o "${APP_NAME}.arm64" -v

}
elseif(param1 -eq 'cleanup') {
    Remove-Item -Path "$APP_NAME.*" -Force
    Remove-Item -Path "dist" -Recurse -Force

}
elseif($args.Length -eq 0) {
    $ret = test_go
    if ($ret -eq 1) {
        exit 1
    }
    protoc
    & $BUILD_COMMAND $BUILD_TAG

} else {
    Write-Output "Unknown argument: ${args[0]}"
    exit 1
}