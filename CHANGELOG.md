# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Planning]
- Configmgr新增從遠端讀取配置(viper的ReadRemoteConfig)
- 第三方登入組件(Apple / Google / FB)
    - Apple
        - https://github.com/pikann/go-verify-apple-id-token/tree/master
        - https://github.com/Timothylock/go-signin-with-apple/tree/master
        - https://blog.csdn.net/weixin_37743259/article/details/123731753
        - https://blog.jks.coffee/sign-in-with-apple/
    - Google
        - https://github.com/googleapis/google-api-go-client/blob/main/oauth2/v2/oauth2-gen.go
        - https://developers.google.com/identity/gsi/web/guides/verify-google-id-token?hl=en
    - FB
        - https://developers.facebook.com/docs/graph-api/reference/v17.0/debug_token
        - https://blog.csdn.net/huqiankunlol/article/details/102721611
- 模組增加current來控制時間

## [Unrelease]
