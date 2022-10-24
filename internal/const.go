package internal

import log "github.com/sirupsen/logrus"

const DefaultRole = "Viewer"
const DefaultAuthSetting = SamlAuthSetting
const DefaultLogLevel = log.WarnLevel

const SamlAuthSetting = "SAML"

const ExistingAssetsUserNameVar = "TABLEAU_EXISTING_ASSETS_USER_NAME"

const EnvVarLogLevel = "LOG_LEVEL"
