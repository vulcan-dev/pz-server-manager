package pz

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"strconv"
)

// Todo for tomorrow: Hardcode the fields with default values, their type and anything else.
// https://ini.unknwon.io/docs/advanced/map_and_reflect
// Will make it easier to read, modify & save the config file.
//PVP                                          bool
//PauseEmpty                                   bool
//GlobalChat                                   bool
//ChatStreams                                  string
//Open                                         bool
//ServerWelcomeMessage                         string
//AutoCreateUserInWhiteList                    bool
//DisplayUserName                              bool
//ShowFirstAndLastName                         bool
//SpawnPoint                                   string
//SafetySystem                                 bool
//SafetyToggleTimer                            int
//SafetyCooldownTimer                          int
//SpawnItems                                   string
//DefaultPort                                  int
//UDPPort                                      int
//ResetID                                      int
//Mods                                         string
//Map                                          string
//DoLuaChecksum                                bool
//DenyLoginOnOverloadedServer                  bool
//Public                                       bool
//PublicName                                   string
//PublicDescription                            string
//MaxPlayers                                   int
//PingLimit                                    int
//HoursForLootRespawn                          int
//MaxItemsForLootRespawn                       int
//ConstructionPreventsLootRespawn              bool
//DropOffWhiteListAfterDeath                   bool
//NoFire                                       bool
//AnnounceDeath                                bool
//MinutesPerPage                               int
//SaveWorldEveryMinutes                        int
//PlayerSafehouse                              bool
//AdminSafehouse                               bool
//SafehouseAllowTrepass                        bool
//SafehouseAllowFire                           bool
//SafehouseAllowLoot                           bool
//SafehouseAllowRespawn                        bool
//SafehouseAllowPVP                            bool
//SafeHouseRemovalTime                         int
//SafehouseAllowNonResidential                 bool
//AllowDestructionBySledgehammer               bool
//SledgehammerOnlyInSafehouse                  bool
//KickFastPlayers                              bool
//ServerPlayerID                               int
//RCONPort                                     int
//RCONPassword                                 string
//DiscordEnable                                bool
//DiscordToken                                 string
//DiscordChannel                               string
//DiscordChannelID                             string
//Password                                     string
//MaxAccountsPerUser                           int
//AllowCoop                                    bool
//SleepAllowed                                 bool
//SleepNeeded                                  bool
//KnockedDownAllowed                           bool
//SneakModeHideFromOtherPlayers                bool
//WorkshopItems                                string
//SteamScoreboard                              bool
//SteamVAC                                     bool
//UPnP                                         bool
//VoiceEnable                                  bool
//VoiceMinDistance                             int
//VoiceMaxDistance                             int
//Voice3D                                      bool
//SpeedLimit                                   int
//LoginQueueEnabled                            bool
//LoginQueueConnectTimeout                     int
//server_browser_announced_ip                  string
//PlayerRespawnWithSelf                        bool
//PlayerRespawnWithOther                       bool
//FastForwardMultiplier                        int
//DisableSafehouseWhenPlayerConnected          bool
//Faction                                      string
//FactionDaySurvivedToCreate                   int
//FactionPlayersRequiredForTag                 int
//DisableRadioStaff                            bool
//DisableRadioAdmin                            bool
//DisableRadioGM                               bool
//DisableRadioOverseer                         bool
//DisableRadioModerator                        bool
//DisableRadioInvisible                        bool
//ClientCommandFilter                          string
//ClientActionLogs                             string
//PerkLogs                                     bool
//ItemNumbersLimitPerContainer                 int
//BloodSplatLifespanDays                       int
//AllowNonAsciiUsername                        bool
//BanKickGlobalSound                           bool
//RemovePlayerCorpsesOnCorpseRemoval           bool
//TrashDeleteAll                               bool
//PVPMeleeWhileHitReaction                     bool
//MouseOverToSeeDisplayName                    bool
//HidePlayersBehindYou                         bool
//PVPMeleeDamageModifier                       float64
//PVPFirearmDamageModifier                     float64
//CarEngineAttractionModifier                  float64
//PlayerBumpPlayer                             bool
//MapRemotePlayerVisibility                    int
//BackupsCount                                 int
//BackupsOnStart                               bool
//BackupsOnVersionChange                       bool
//BackupsPeriod                                int
//AntiCheatProtectionType1                     bool
//AntiCheatProtectionType2                     bool
//AntiCheatProtectionType3                     bool
//AntiCheatProtectionType4                     bool
//AntiCheatProtectionType5                     bool
//AntiCheatProtectionType6                     bool
//AntiCheatProtectionType7                     bool
//AntiCheatProtectionType8                     bool
//AntiCheatProtectionType9                     bool
//AntiCheatProtectionType10                    bool
//AntiCheatProtectionType11                    bool
//AntiCheatProtectionType12                    bool
//AntiCheatProtectionType13                    bool
//AntiCheatProtectionType14                    bool
//AntiCheatProtectionType15                    bool
//AntiCheatProtectionType16                    bool
//AntiCheatProtectionType17                    bool
//AntiCheatProtectionType18                    bool
//AntiCheatProtectionType19                    bool
//AntiCheatProtectionType20                    bool
//AntiCheatProtectionType21                    bool
//AntiCheatProtectionType22                    bool
//AntiCheatProtectionType23                    bool
//AntiCheatProtectionType24                    bool
//AntiCheatProtectionType2ThresholdMultiplier  float64
//AntiCheatProtectionType3ThresholdMultiplier  float64
//AntiCheatProtectionType4ThresholdMultiplier  float64
//AntiCheatProtectionType9ThresholdMultiplier  float64
//AntiCheatProtectionType15ThresholdMultiplier float64
//AntiCheatProtectionType20ThresholdMultiplier float64
//AntiCheatProtectionType22ThresholdMultiplier float64
//AntiCheatProtectionType24ThresholdMultiplier float64

import (
	"github.com/go-ini/ini"
)

type Variable struct {
	Name  string
	Value interface{}
}

type Config struct {
	config map[string]Variable
}

func Parse(name string, path string) (Config, error) {
	path = path + "/Server/" + name + ".ini"
	iniFile, err := ini.Load(path)
	if err != nil {
		return Config{}, err
	}

	config := Config{}
	config.config = make(map[string]Variable)

	for _, section := range iniFile.Sections() {
		for _, key := range section.Keys() {
			val := key.Value()
			valType := "string"
			if val == "true" || val == "false" {
				valType = "bool"
			} else if _, err := strconv.Atoi(val); err == nil {
				valType = "int"
			} else if _, err := strconv.ParseFloat(val, 64); err == nil {
				valType = "float"
			} else if _, err := strconv.ParseBool(val); err == nil {
				valType = "bool"
			}

			switch valType {
			case "int":
				intVal, _ := strconv.Atoi(val)
				config.config[key.Name()] = Variable{key.Name(), intVal}
			case "float":
				floatVal, _ := strconv.ParseFloat(val, 64)
				config.config[key.Name()] = Variable{key.Name(), floatVal}
			case "bool":
				boolVal, _ := strconv.ParseBool(val)
				config.config[key.Name()] = Variable{key.Name(), boolVal}
			default:
				config.config[key.Name()] = Variable{key.Name(), val}
			}
		}
	}

	fmt.Println(config.ReadToJSONString())

	return config, nil
}

func (c Config) ReadToJSONString() string {
	jsonString := "{\n"
	for _, v := range c.config {
		_, ok := v.Value.(string)
		if ok {
			v.Value = string([]rune(v.Value.(string)))
		}

		jsonString += fmt.Sprintf("\"%s\": \"%v\",\n", v.Name, v.Value)
	}
	jsonString = jsonString[:len(jsonString)-2]
	jsonString += "\n}"
	return jsonString
}

func (c Config) ReadToGIN(ginContext *gin.Context) map[string]interface{} {
	json := make(map[string]interface{})
	for _, v := range c.config {
		_, ok := v.Value.(string)
		if ok {
			v.Value = string([]rune(v.Value.(string)))
		}

		json[v.Name] = v.Value
	}
	return json
}

func (c Config) Get(name string) Variable {
	return c.config[name]
}

func (c Config) GetInt(name string) int {
	return c.config[name].Value.(int)
}

func (c Config) GetFloat(name string) float64 {
	return c.config[name].Value.(float64)
}

func (c Config) GetBool(name string) bool {
	return c.config[name].Value.(bool)
}

func (c Config) WriteString(name string, value string) error {
	if value == "" {
		return nil
	}

	_, ok := c.config[name]
	if !ok {
		return fmt.Errorf("variable %s does not exist", name)
	}

	c.config[name] = Variable{name, value}
	return nil
}

func (c Config) WriteInt(name string, value int) error {
	_, ok := c.config[name]
	if !ok {
		return fmt.Errorf("variable %s does not exist", name)
	}

	c.config[name] = Variable{name, value}
	return nil
}

func (c Config) WriteFloat(name string, value float64) error {
	_, ok := c.config[name]
	if !ok {
		return fmt.Errorf("variable %s does not exist", name)
	}

	c.config[name] = Variable{name, value}
	return nil
}

func (c Config) WriteBool(name string, value bool) error {
	_, ok := c.config[name]
	if !ok {
		return fmt.Errorf("variable %s does not exist", name)
	}

	c.config[name] = Variable{name, value}
	return nil
}

func (c Config) Save(name string, path string) error {
	path = path + "/Server/" + name + ".ini"
	iniFile, err := ini.Load(path)
	if err != nil {
		return err
	}

	for _, section := range iniFile.Sections() {
		for _, key := range section.Keys() {
			val := c.Get(key.Name())
			section.Key(key.Name()).SetValue(fmt.Sprintf("%v", val.Value))
		}
	}

	return iniFile.SaveTo(path)
}

func (c Config) SaveFromJSON(name string, path string, json map[string]interface{}) error {
	for key, value := range json {
		log.Debugf("%s = %v\n", key, value)
		valType := "string"
		if value == "true" || value == "false" {
			valType = "bool"
		} else if _, err := strconv.Atoi(value.(string)); err == nil {
			valType = "int"
		} else if _, err := strconv.ParseFloat(value.(string), 64); err == nil {
			valType = "float"
		} else if _, err := strconv.ParseBool(value.(string)); err == nil {
			valType = "bool"
		}

		switch valType {
		case "int":
			intVal, _ := strconv.Atoi(value.(string))
			c.config[key] = Variable{key, intVal}
		case "float":
			floatVal, _ := strconv.ParseFloat(value.(string), 64)
			c.config[key] = Variable{key, floatVal}
		case "bool":
			boolVal, _ := strconv.ParseBool(value.(string))
			c.config[key] = Variable{key, boolVal}
		default:
			c.config[key] = Variable{key, value.(string)}
		}
	}

	return c.Save(name, path)
}
