# Mail Sender API

This is simple API for sending email notifications via predefined user configured during the deployment.

# Server configuration:

- The password secret contains the username and the password of the user from whom the email notification will be send. Currently only plain vault for retreving the secrets is used. The server is with configured TLS. The certificates are expected to be under `./cfg/tls/sec`

    - Server configuration example:
    ```
    vaultType: plain
    SMTPServerCfg:
    host: smtp.gmail.com
    port: 587
    passwordSecret: passwordSecret
    keepAlive: true
    connectionTimeout: 30s
    sendTimeout: 30s
    security:
    #If tls configuration is not provided HTTP server communication is configured
    tls:
        #The cert file is read from ./cfg/tls directly.
        #If the certificate is signed by a certificate authority,
        #the certFile should be the concatenation of the server's certificate,
        #any intermediates, and the CA's certificate.
        certFile: mailSenderCert.pem
        #The matching for the certificate private key.
        privateKey: mailSenderKey.pem
        rootCACert: rootCACert.pem
        rootCAKey: rootCAKey.pem
    ```

    - Vault configuration:
     ```
        - id: <secretID>
          name: <name>
          value: <value>
        ```

# Docker image

- Docker image could be found [here](https://hub.docker.com/repository/docker/todorov99/mailsender)