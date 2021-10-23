package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/MohammedAl-Mahdawi/bnkr/app/dal"
	"github.com/MohammedAl-Mahdawi/bnkr/app/types"
	"github.com/MohammedAl-Mahdawi/bnkr/config"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

var app *config.AppConfig

// NewUtils sets up app config for helpers
func NewUtils(a *config.AppConfig) {
	app = a
}

func GetTimes() []string {
	t := []string{
		"00:00",
		"00:30",
		"01:00",
		"01:30",
		"02:00",
		"02:30",
		"03:00",
		"03:30",
		"04:00",
		"04:30",
		"05:00",
		"05:30",
		"06:00",
		"06:30",
		"07:00",
		"07:30",
		"08:00",
		"08:30",
		"09:00",
		"09:30",
		"10:00",
		"10:30",
		"11:00",
		"11:30",
		"12:00",
		"12:30",
		"13:00",
		"13:30",
		"14:00",
		"14:30",
		"15:00",
		"15:30",
		"16:00",
		"16:30",
		"17:00",
		"17:30",
		"18:00",
		"18:30",
		"19:00",
		"19:30",
		"20:00",
		"20:30",
		"21:00",
		"21:30",
		"22:00",
		"22:30",
		"23:00",
		"23:30",
	}

	return t
}

func GetTimeZones() []string {
	timezones := []string{
		"Africa/Abidjan",
		"Africa/Accra",
		"Africa/Addis_Ababa",
		"Africa/Algiers",
		"Africa/Asmara",
		"Africa/Bamako",
		"Africa/Bangui",
		"Africa/Banjul",
		"Africa/Bissau",
		"Africa/Blantyre",
		"Africa/Brazzaville",
		"Africa/Bujumbura",
		"Africa/Cairo",
		"Africa/Casablanca",
		"Africa/Ceuta",
		"Africa/Conakry",
		"Africa/Dakar",
		"Africa/Dar_es_Salaam",
		"Africa/Djibouti",
		"Africa/Douala",
		"Africa/El_Aaiun",
		"Africa/Freetown",
		"Africa/Gaborone",
		"Africa/Harare",
		"Africa/Johannesburg",
		"Africa/Juba",
		"Africa/Kampala",
		"Africa/Khartoum",
		"Africa/Kigali",
		"Africa/Kinshasa",
		"Africa/Lagos",
		"Africa/Libreville",
		"Africa/Lome",
		"Africa/Luanda",
		"Africa/Lubumbashi",
		"Africa/Lusaka",
		"Africa/Malabo",
		"Africa/Maputo",
		"Africa/Maseru",
		"Africa/Mbabane",
		"Africa/Mogadishu",
		"Africa/Monrovia",
		"Africa/Nairobi",
		"Africa/Ndjamena",
		"Africa/Niamey",
		"Africa/Nouakchott",
		"Africa/Ouagadougou",
		"Africa/Porto-Novo",
		"Africa/Sao_Tome",
		"Africa/Tripoli",
		"Africa/Tunis",
		"Africa/Windhoek",
		"America/Adak",
		"America/Anchorage",
		"America/Anguilla",
		"America/Antigua",
		"America/Araguaina",
		"America/Argentina/Buenos_Aires",
		"America/Argentina/Catamarca",
		"America/Argentina/Cordoba",
		"America/Argentina/Jujuy",
		"America/Argentina/La_Rioja",
		"America/Argentina/Mendoza",
		"America/Argentina/Rio_Gallegos",
		"America/Argentina/Salta",
		"America/Argentina/San_Juan",
		"America/Argentina/San_Luis",
		"America/Argentina/Tucuman",
		"America/Argentina/Ushuaia",
		"America/Aruba",
		"America/Asuncion",
		"America/Atikokan",
		"America/Bahia",
		"America/Bahia_Banderas",
		"America/Barbados",
		"America/Belem",
		"America/Belize",
		"America/Blanc-Sablon",
		"America/Boa_Vista",
		"America/Bogota",
		"America/Boise",
		"America/Cambridge_Bay",
		"America/Campo_Grande",
		"America/Cancun",
		"America/Caracas",
		"America/Cayenne",
		"America/Cayman",
		"America/Chicago",
		"America/Chihuahua",
		"America/Costa_Rica",
		"America/Creston",
		"America/Cuiaba",
		"America/Curacao",
		"America/Danmarkshavn",
		"America/Dawson",
		"America/Dawson_Creek",
		"America/Denver",
		"America/Detroit",
		"America/Dominica",
		"America/Edmonton",
		"America/Eirunepe",
		"America/El_Salvador",
		"America/Fort_Nelson",
		"America/Fortaleza",
		"America/Glace_Bay",
		"America/Goose_Bay",
		"America/Grand_Turk",
		"America/Grenada",
		"America/Guadeloupe",
		"America/Guatemala",
		"America/Guayaquil",
		"America/Guyana",
		"America/Halifax",
		"America/Havana",
		"America/Hermosillo",
		"America/Indiana/Indianapolis",
		"America/Indiana/Knox",
		"America/Indiana/Marengo",
		"America/Indiana/Petersburg",
		"America/Indiana/Tell_City",
		"America/Indiana/Vevay",
		"America/Indiana/Vincennes",
		"America/Indiana/Winamac",
		"America/Inuvik",
		"America/Iqaluit",
		"America/Jamaica",
		"America/Juneau",
		"America/Kentucky/Louisville",
		"America/Kentucky/Monticello",
		"America/Kralendijk",
		"America/La_Paz",
		"America/Lima",
		"America/Los_Angeles",
		"America/Lower_Princes",
		"America/Maceio",
		"America/Managua",
		"America/Manaus",
		"America/Marigot",
		"America/Martinique",
		"America/Matamoros",
		"America/Mazatlan",
		"America/Menominee",
		"America/Merida",
		"America/Metlakatla",
		"America/Mexico_City",
		"America/Miquelon",
		"America/Moncton",
		"America/Monterrey",
		"America/Montevideo",
		"America/Montserrat",
		"America/Nassau",
		"America/New_York",
		"America/Nipigon",
		"America/Nome",
		"America/Noronha",
		"America/North_Dakota/Beulah",
		"America/North_Dakota/Center",
		"America/North_Dakota/New_Salem",
		"America/Nuuk",
		"America/Ojinaga",
		"America/Panama",
		"America/Pangnirtung",
		"America/Paramaribo",
		"America/Phoenix",
		"America/Port-au-Prince",
		"America/Port_of_Spain",
		"America/Porto_Velho",
		"America/Puerto_Rico",
		"America/Punta_Arenas",
		"America/Rainy_River",
		"America/Rankin_Inlet",
		"America/Recife",
		"America/Regina",
		"America/Resolute",
		"America/Rio_Branco",
		"America/Santarem",
		"America/Santiago",
		"America/Santo_Domingo",
		"America/Sao_Paulo",
		"America/Scoresbysund",
		"America/Sitka",
		"America/St_Barthelemy",
		"America/St_Johns",
		"America/St_Kitts",
		"America/St_Lucia",
		"America/St_Thomas",
		"America/St_Vincent",
		"America/Swift_Current",
		"America/Tegucigalpa",
		"America/Thule",
		"America/Thunder_Bay",
		"America/Tijuana",
		"America/Toronto",
		"America/Tortola",
		"America/Vancouver",
		"America/Whitehorse",
		"America/Winnipeg",
		"America/Yakutat",
		"America/Yellowknife",
		"Antarctica/Casey",
		"Antarctica/Davis",
		"Antarctica/DumontDUrville",
		"Antarctica/Macquarie",
		"Antarctica/Mawson",
		"Antarctica/McMurdo",
		"Antarctica/Palmer",
		"Antarctica/Rothera",
		"Antarctica/Syowa",
		"Antarctica/Troll",
		"Antarctica/Vostok",
		"Arctic/Longyearbyen",
		"Asia/Aden",
		"Asia/Almaty",
		"Asia/Amman",
		"Asia/Anadyr",
		"Asia/Aqtau",
		"Asia/Aqtobe",
		"Asia/Ashgabat",
		"Asia/Atyrau",
		"Asia/Baghdad",
		"Asia/Bahrain",
		"Asia/Baku",
		"Asia/Bangkok",
		"Asia/Barnaul",
		"Asia/Beirut",
		"Asia/Bishkek",
		"Asia/Brunei",
		"Asia/Chita",
		"Asia/Choibalsan",
		"Asia/Colombo",
		"Asia/Damascus",
		"Asia/Dhaka",
		"Asia/Dili",
		"Asia/Dubai",
		"Asia/Dushanbe",
		"Asia/Famagusta",
		"Asia/Gaza",
		"Asia/Hebron",
		"Asia/Ho_Chi_Minh",
		"Asia/Hong_Kong",
		"Asia/Hovd",
		"Asia/Irkutsk",
		"Asia/Jakarta",
		"Asia/Jayapura",
		"Asia/Jerusalem",
		"Asia/Kabul",
		"Asia/Kamchatka",
		"Asia/Karachi",
		"Asia/Kathmandu",
		"Asia/Khandyga",
		"Asia/Kolkata",
		"Asia/Krasnoyarsk",
		"Asia/Kuala_Lumpur",
		"Asia/Kuching",
		"Asia/Kuwait",
		"Asia/Macau",
		"Asia/Magadan",
		"Asia/Makassar",
		"Asia/Manila",
		"Asia/Muscat",
		"Asia/Nicosia",
		"Asia/Novokuznetsk",
		"Asia/Novosibirsk",
		"Asia/Omsk",
		"Asia/Oral",
		"Asia/Phnom_Penh",
		"Asia/Pontianak",
		"Asia/Pyongyang",
		"Asia/Qatar",
		"Asia/Qostanay",
		"Asia/Qyzylorda",
		"Asia/Riyadh",
		"Asia/Sakhalin",
		"Asia/Samarkand",
		"Asia/Seoul",
		"Asia/Shanghai",
		"Asia/Singapore",
		"Asia/Srednekolymsk",
		"Asia/Taipei",
		"Asia/Tashkent",
		"Asia/Tbilisi",
		"Asia/Tehran",
		"Asia/Thimphu",
		"Asia/Tokyo",
		"Asia/Tomsk",
		"Asia/Ulaanbaatar",
		"Asia/Urumqi",
		"Asia/Ust-Nera",
		"Asia/Vientiane",
		"Asia/Vladivostok",
		"Asia/Yakutsk",
		"Asia/Yangon",
		"Asia/Yekaterinburg",
		"Asia/Yerevan",
		"Atlantic/Azores",
		"Atlantic/Bermuda",
		"Atlantic/Canary",
		"Atlantic/Cape_Verde",
		"Atlantic/Faroe",
		"Atlantic/Madeira",
		"Atlantic/Reykjavik",
		"Atlantic/South_Georgia",
		"Atlantic/St_Helena",
		"Atlantic/Stanley",
		"Australia/Adelaide",
		"Australia/Brisbane",
		"Australia/Broken_Hill",
		"Australia/Darwin",
		"Australia/Eucla",
		"Australia/Hobart",
		"Australia/Lindeman",
		"Australia/Lord_Howe",
		"Australia/Melbourne",
		"Australia/Perth",
		"Australia/Sydney",
		"Europe/Amsterdam",
		"Europe/Andorra",
		"Europe/Astrakhan",
		"Europe/Athens",
		"Europe/Belgrade",
		"Europe/Berlin",
		"Europe/Bratislava",
		"Europe/Brussels",
		"Europe/Bucharest",
		"Europe/Budapest",
		"Europe/Busingen",
		"Europe/Chisinau",
		"Europe/Copenhagen",
		"Europe/Dublin",
		"Europe/Gibraltar",
		"Europe/Guernsey",
		"Europe/Helsinki",
		"Europe/Isle_of_Man",
		"Europe/Istanbul",
		"Europe/Jersey",
		"Europe/Kaliningrad",
		"Europe/Kiev",
		"Europe/Kirov",
		"Europe/Lisbon",
		"Europe/Ljubljana",
		"Europe/London",
		"Europe/Luxembourg",
		"Europe/Madrid",
		"Europe/Malta",
		"Europe/Mariehamn",
		"Europe/Minsk",
		"Europe/Monaco",
		"Europe/Moscow",
		"Europe/Oslo",
		"Europe/Paris",
		"Europe/Podgorica",
		"Europe/Prague",
		"Europe/Riga",
		"Europe/Rome",
		"Europe/Samara",
		"Europe/San_Marino",
		"Europe/Sarajevo",
		"Europe/Saratov",
		"Europe/Simferopol",
		"Europe/Skopje",
		"Europe/Sofia",
		"Europe/Stockholm",
		"Europe/Tallinn",
		"Europe/Tirane",
		"Europe/Ulyanovsk",
		"Europe/Uzhgorod",
		"Europe/Vaduz",
		"Europe/Vatican",
		"Europe/Vienna",
		"Europe/Vilnius",
		"Europe/Volgograd",
		"Europe/Warsaw",
		"Europe/Zagreb",
		"Europe/Zaporozhye",
		"Europe/Zurich",
		"Indian/Antananarivo",
		"Indian/Chagos",
		"Indian/Christmas",
		"Indian/Cocos",
		"Indian/Comoro",
		"Indian/Kerguelen",
		"Indian/Mahe",
		"Indian/Maldives",
		"Indian/Mauritius",
		"Indian/Mayotte",
		"Indian/Reunion",
		"Pacific/Apia",
		"Pacific/Auckland",
		"Pacific/Bougainville",
		"Pacific/Chatham",
		"Pacific/Chuuk",
		"Pacific/Easter",
		"Pacific/Efate",
		"Pacific/Enderbury",
		"Pacific/Fakaofo",
		"Pacific/Fiji",
		"Pacific/Funafuti",
		"Pacific/Galapagos",
		"Pacific/Gambier",
		"Pacific/Guadalcanal",
		"Pacific/Guam",
		"Pacific/Honolulu",
		"Pacific/Kiritimati",
		"Pacific/Kosrae",
		"Pacific/Kwajalein",
		"Pacific/Majuro",
		"Pacific/Marquesas",
		"Pacific/Midway",
		"Pacific/Nauru",
		"Pacific/Niue",
		"Pacific/Norfolk",
		"Pacific/Noumea",
		"Pacific/Pago_Pago",
		"Pacific/Palau",
		"Pacific/Pitcairn",
		"Pacific/Pohnpei",
		"Pacific/Port_Moresby",
		"Pacific/Rarotonga",
		"Pacific/Saipan",
		"Pacific/Tahiti",
		"Pacific/Tarawa",
		"Pacific/Tongatapu",
		"Pacific/Wake",
		"Pacific/Wallis",
		"UTC",
	}

	return timezones
}

// GetUser is helper function for getting authenticated user's id
func GetUser(w http.ResponseWriter, r *http.Request) int {
	userId, ok := app.Session.Get(r.Context(), "user_id").(int)
	if !ok {
		app.Session.Put(r.Context(), "error", "Can't get user id from session!")
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return 0
	}
	return userId
}

func IsAuthenticated(r *http.Request) bool {
	exists := app.Session.Exists(r.Context(), "user_id")
	return exists
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func IsDuplicateKeyError(err error) bool {
	duplicate := regexp.MustCompile(`unique constraint`)
	return duplicate.MatchString(err.Error())
}

func GetOptionValue(o string) string {
	option := &types.NewOptionDTO{}
	if err := dal.FindOptionByName(option, o); err != nil {
		return ""
	}

	return option.Value
}

func GetRequiredMigAccessFields(values *dal.Migration) []string {
	result := []string{}

	if values.SrcType == "db" || values.SrcType == "pg" || values.SrcType == "mongo" || values.SrcType == "bnkr" {
		if values.SrcAccess == "k8s" {
			result = []string{"src_kubeconfig"}
		}

		if values.SrcAccess == "ssh" {
			result = []string{"src_ssh_host", "src_ssh_port", "src_ssh_user", "src_ssh_key"}
		}
	}

	if values.DestType == "db" || values.DestType == "pg" || values.DestType == "mongo" || values.DestType == "bnkr" {
		if values.DestAccess == "k8s" {
			result = append(result, "dest_kubeconfig")
		}

		if values.DestAccess == "ssh" {
			result = append(result, "dest_ssh_host", "dest_ssh_port", "dest_ssh_user", "dest_ssh_key")
		}
	}

	return result
}

func GetRequiredMigTypeFields(theType string, itfor string) []string {
	result := []string{}

	switch theType {
	case "db":
		result = []string{itfor + "_db_name", itfor + "_db_user", itfor + "_db_password", itfor + "_db_host", itfor + "_db_port"}
	case "object":
		result = []string{itfor + "_pod_label", itfor + "_files_path", itfor + "_container"}
	case "mongo", "pg":
		result = []string{itfor + "_uri"}
	case "pod":
		result = []string{itfor + "_files_path", itfor + "_container", itfor + "_pod_name"}
	case "ssh":
		result = []string{itfor + "_ssh_host", itfor + "_ssh_port", itfor + "_ssh_user", itfor + "_ssh_key"}
	case "s3":
		result = []string{itfor + "_bucket", itfor + "_s3_access_key", itfor + "_s3_secret_key", itfor + "_region"}
		if itfor == "src" {
			result = append(result, itfor+"_s3_file")
		}
	}

	return result
}

func CreateKubeconfigFile(path string, kubeconfig string, kubeconfigName string) (string, error) {
	if kubeconfigName == "" {
		kubeconfigName = "kubeconfig.yml"
	}
	kubeconfigPath := path + "/" + kubeconfigName
	err := ioutil.WriteFile(kubeconfigPath, []byte(kubeconfig), 0600)
	if err != nil {
		return "", err
	}

	return kubeconfigPath, nil
}

func CreateSSHKeyFile(path string, key string, keyName string) (string, error) {
	if keyName == "" {
		keyName = "id_rsa"
	}

	sshKeyPath := path + "/" + keyName
	err := ioutil.WriteFile(sshKeyPath, []byte(key), 0600)
	if err != nil {
		return "", err
	}

	return sshKeyPath, nil
}

func CmdExecutor(cmd *exec.Cmd) (string, error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		app.ErrorLog.Println(err)
		return err.Error(), err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		app.ErrorLog.Println(err)
		return err.Error(), err
	}

	if err := cmd.Start(); err != nil {
		app.ErrorLog.Println(err)
		return err.Error(), err
	}

	s := bufio.NewScanner(io.MultiReader(stdout, stderr))
	var co bytes.Buffer

	for s.Scan() {
		co.Write(s.Bytes())
		co.Write([]byte("\n"))
	}

	if err := cmd.Wait(); err != nil {
		cerr := co.String() + `

		` + err.Error()

		return cerr, err
	}

	return co.String(), err
}

func ErrorJSON(w http.ResponseWriter, err error) {
	type jsonError struct {
		Message string `json:"message"`
	}

	theError := jsonError{
		Message: err.Error(),
	}

	WriteJSON(w, http.StatusBadRequest, theError, "error")
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}, wrap string) error {
	wrapper := make(map[string]interface{})

	wrapper[wrap] = data

	js, err := json.Marshal(wrapper)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func GetSshConn(user string, host string, port string, key string) (*ssh.Client, error) {
	signer, err := ssh.ParsePrivateKey([]byte(key))
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", host+":"+port, config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func RunSshCommand(conn *ssh.Client, cmd string) (string, error) {
	var o string
	buf := new(strings.Builder)

	sess, err := conn.NewSession()
	if err != nil {
		return o, err
	}
	defer sess.Close()
	sessStdOut, err := sess.StdoutPipe()
	if err != nil {
		return o, err
	}

	sessStderr, err := sess.StderrPipe()
	if err != nil {
		return o, err
	}

	r := io.MultiReader(sessStdOut, sessStderr)
	_, err = io.Copy(buf, r)
	if err != nil {
		return o, err
	}
	o = buf.String()

	err = sess.Run(cmd)
	if err != nil {
		return o, err
	}

	return o, nil
}

func UploadFileOverSftp(conn *ssh.Client, sFile string, dFile string) (int64, error) {
	client, err := sftp.NewClient(conn)
	if err != nil {
		return 0, err
	}
	defer client.Close()

	// create destination file
	dstFile, err := client.Create(dFile)
	if err != nil {
		return 0, err
	}
	defer dstFile.Close()

	// open source file
	srcFile, err := os.Open(sFile)
	if err != nil {
		return 0, err
	}

	// copy source file to destination file
	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		return 0, err
	}

	return bytes, nil
}

func DownloadFileOverSftp(conn *ssh.Client, sFile string, dFile string) (int64, error) {
	client, err := sftp.NewClient(conn)
	if err != nil {
		return 0, err
	}
	defer client.Close()

	// create destination file
	dstFile, err := os.Create(dFile)
	if err != nil {
		return 0, err
	}
	defer dstFile.Close()

	// open source file
	srcFile, err := client.Open(sFile)
	if err != nil {
		return 0, err
	}

	// copy source file to destination file
	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		return 0, err
	}

	// flush in-memory copy
	err = dstFile.Sync()
	if err != nil {
		return bytes, err
	}

	return bytes, nil
}
