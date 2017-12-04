package ldap_test

import (
	"io/ioutil"
	"os"
	"strconv"

	. "github.com/pivotalservices/cf-mgmt/ldap"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ldap", func() {
	var ldapManager Manager
	var config *Config
	Describe("given a ldap manager", func() {
		BeforeEach(func() {
			var host string
			var port int
			if os.Getenv("LDAP_PORT_389_TCP_ADDR") == "" {
				host = "127.0.0.1"
				port = 389
			} else {
				host = os.Getenv("LDAP_PORT_389_TCP_ADDR")
				port, _ = strconv.Atoi(os.Getenv("LDAP_PORT_389_TCP_PORT"))
			}
			ldapManager = &DefaultManager{}
			config = &Config{
				BindDN:            "cn=admin,dc=pivotal,dc=org",
				BindPassword:      "password",
				UserSearchBase:    "dc=pivotal,dc=org",
				UserNameAttribute: "uid",
				UserMailAttribute: "mail",
				GroupSearchBase:   "ou=groups,dc=pivotal,dc=org",
				GroupAttribute:    "member",
				LdapHost:          host,
				LdapPort:          port,
			}
		})
		Context("when cn with special characters", func() {
			It("then it should return 1 Entry", func() {
				entry, err := ldapManager.GetLdapUser(config, "cn=Washburn, Caleb,ou=users,dc=pivotal,dc=org")
				Ω(err).Should(BeNil())
				Ω(entry).ShouldNot(BeNil())
			})
		})
		Context("when cn has a period", func() {
			It("then it should return 1 Entry", func() {
				entry, err := ldapManager.GetLdapUser(config, "cn=Caleb A. Washburn,ou=users,dc=pivotal,dc=org")
				Ω(err).Should(BeNil())
				Ω(entry).ShouldNot(BeNil())
			})
		})
		Context("when called with a valid group", func() {
			It("then it should return 5 users", func() {
				users, err := ldapManager.GetUserIDs(config, "space_developers")
				Ω(err).Should(BeNil())
				Ω(len(users)).Should(Equal(5))
			})
		})
		Context("when called with a valid group with special characters", func() {
			It("then it should return 4 users", func() {
				users, err := ldapManager.GetUserIDs(config, "special (char) group,name")
				Ω(err).Should(BeNil())
				Ω(len(users)).Should(Equal(4))
			})
		})
		Context("GetUser()", func() {
			It("then it should return 1 user", func() {
				user, err := ldapManager.GetUser(config, "cwashburn")
				Ω(err).Should(BeNil())
				Ω(user).ShouldNot(BeNil())
				Ω(user.UserID).Should(Equal("cwashburn"))
				Ω(user.UserDN).Should(Equal("cn=cwashburn,ou=users,dc=pivotal,dc=org"))
				Ω(user.Email).Should(Equal("cwashburn+cfmt@testdomain.com"))
			})
		})

		Describe("given a ldap manager with userObjectClass", func() {
			BeforeEach(func() {
				var host string
				var port int
				if os.Getenv("LDAP_PORT_389_TCP_ADDR") == "" {
					host = "127.0.0.1"
					port = 389
				} else {
					host = os.Getenv("LDAP_PORT_389_TCP_ADDR")
					port, _ = strconv.Atoi(os.Getenv("LDAP_PORT_389_TCP_PORT"))
				}
				ldapManager = &DefaultManager{}
				config = &Config{
					BindDN:            "cn=admin,dc=pivotal,dc=org",
					BindPassword:      "password",
					UserSearchBase:    "dc=pivotal,dc=org",
					UserNameAttribute: "uid",
					UserMailAttribute: "mail",
					GroupSearchBase:   "ou=groups,dc=pivotal,dc=org",
					GroupAttribute:    "member",
					LdapHost:          host,
					LdapPort:          port,
					UserObjectClass:   "inetOrgPerson",
				}
			})
			Context("when cn with special characters", func() {
				It("then it should return 1 Entry", func() {
					entry, err := ldapManager.GetLdapUser(config, "cn=Washburn, Caleb,ou=users,dc=pivotal,dc=org")
					Ω(err).Should(BeNil())
					Ω(entry).ShouldNot(BeNil())
				})
			})
			Context("GetUser()", func() {
				It("then it should return 1 user", func() {
					user, err := ldapManager.GetUser(config, "cwashburn")
					Ω(err).Should(BeNil())
					Ω(user).ShouldNot(BeNil())
					Ω(user.UserID).Should(Equal("cwashburn"))
					Ω(user.UserDN).Should(Equal("cn=cwashburn,ou=users,dc=pivotal,dc=org"))
					Ω(user.Email).Should(Equal("cwashburn+cfmt@testdomain.com"))
				})
			})
		})
		Context("GetLdapUser()", func() {
			It("then it should return 1 user", func() {
				data, _ := ioutil.ReadFile("./fixtures/user1.txt")
				user, err := ldapManager.GetLdapUser(config, string(data))
				Ω(err).Should(BeNil())
				Ω(user).ShouldNot(BeNil())
				Ω(user.UserID).Should(Equal("cwashburn2"))
				Ω(user.Email).Should(Equal("cwashburn+cfmt2@testdomain.com"))
			})
		})
	})
})
