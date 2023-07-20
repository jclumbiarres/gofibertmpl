package config

// Esto tiene que ir en un .env luego cargarlo y parsearlo, prueba de concepto de conexión de base de datos.
const PostgreHOST = "localhost"
const PostgreUSER = "usuario"
const PostgrePASS = "contraseña"
const PostgreDBNAME = "basededatos"
const PostgrePORT = "5432"

const PostgreCONNTXT = "host=" + PostgreHOST + " user=" + PostgreUSER + " password=" + PostgrePASS + " dbname=" + PostgreDBNAME + " port=" + PostgrePORT + " sslmode=disable"
