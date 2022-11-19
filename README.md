# screenpng (simple http server with save screenshot in png)

## This application ([spf13/cobra](https://github.com/spf13/cobra) + [pf13/viper](https://github.com/spf13/viper)) is a simple example of using a [pog7x/ssfactory](https://github.com/pog7x/ssfactory)

### Example of build/run app in docker container ([base docker image](https://hub.docker.com/repository/docker/pog7x/gobasebrowser))

```bash
docker build -t screen_png .

docker run -v $(pwd):/screenpng -p 8099:8099 -it screen_png

```