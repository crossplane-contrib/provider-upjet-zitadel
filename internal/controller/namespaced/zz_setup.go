// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	executionevent "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/action/executionevent"
	executionfunction "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/action/executionfunction"
	executionrequest "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/action/executionrequest"
	executionresponse "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/action/executionresponse"
	target "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/action/target"
	targetpublickey "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/action/targetpublickey"
	webkey "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/active/webkey"
	api "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/application/api"
	key "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/application/key"
	oidc "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/application/oidc"
	saml "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/application/saml"
	domainclaimedmessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/default/domainclaimedmessagetext"
	domainpolicy "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/default/domainpolicy"
	initmessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/default/initmessagetext"
	inviteusermessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/default/inviteusermessagetext"
	labelpolicy "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/default/labelpolicy"
	lockoutpolicy "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/default/lockoutpolicy"
	loginpolicy "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/default/loginpolicy"
	logintexts "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/default/logintexts"
	notificationpolicy "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/default/notificationpolicy"
	oidcsettings "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/default/oidcsettings"
	passwordagepolicy "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/default/passwordagepolicy"
	passwordchangemessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/default/passwordchangemessagetext"
	passwordcomplexitypolicy "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/default/passwordcomplexitypolicy"
	passwordlessregistrationmessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/default/passwordlessregistrationmessagetext"
	passwordresetmessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/default/passwordresetmessagetext"
	privacypolicy "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/default/privacypolicy"
	securitysettings "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/default/securitysettings"
	verifyemailmessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/default/verifyemailmessagetext"
	verifyemailotpmessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/default/verifyemailotpmessagetext"
	verifyphonemessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/default/verifyphonemessagetext"
	verifysmsotpmessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/default/verifysmsotpmessagetext"
	claimedmessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/domain/claimedmessagetext"
	policy "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/domain/policy"
	providerhttp "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/email/providerhttp"
	providersmtp "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/email/providersmtp"
	apple "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/idp/apple"
	azuread "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/idp/azuread"
	github "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/idp/github"
	githubes "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/idp/githubes"
	gitlab "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/idp/gitlab"
	gitlabselfhosted "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/idp/gitlabselfhosted"
	google "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/idp/google"
	ldap "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/idp/ldap"
	oauth "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/idp/oauth"
	oidcidp "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/idp/oidc"
	samlidp "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/idp/saml"
	messagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/init/messagetext"
	customdomain "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/instance/customdomain"
	features "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/instance/features"
	member "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/instance/member"
	restrictions "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/instance/restrictions"
	secretgenerator "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/instance/secretgenerator"
	trusteddomain "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/instance/trusteddomain"
	policylabel "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/label/policy"
	policylockout "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/lockout/policy"
	policylogin "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/login/policy"
	texts "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/login/texts"
	policynotification "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/notification/policy"
	idpapple "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/org/idpapple"
	idpazuread "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/org/idpazuread"
	idpgithub "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/org/idpgithub"
	idpgithubes "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/org/idpgithubes"
	idpgitlab "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/org/idpgitlab"
	idpgitlabselfhosted "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/org/idpgitlabselfhosted"
	idpgoogle "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/org/idpgoogle"
	idpjwt "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/org/idpjwt"
	idpldap "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/org/idpldap"
	idpoauth "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/org/idpoauth"
	idpoidc "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/org/idpoidc"
	idpsaml "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/org/idpsaml"
	memberorg "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/org/member"
	metadata "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/org/metadata"
	org "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/org/org"
	organization "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/org/organization"
	domain "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/organization/domain"
	metadataorganization "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/organization/metadata"
	agepolicy "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/password/agepolicy"
	changemessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/password/changemessagetext"
	complexitypolicy "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/password/complexitypolicy"
	resetmessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/password/resetmessagetext"
	registrationmessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/passwordless/registrationmessagetext"
	policyprivacy "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/privacy/policy"
	grant "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/project/grant"
	grantmember "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/project/grantmember"
	memberproject "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/project/member"
	project "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/project/project"
	role "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/project/role"
	providerconfig "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/providerconfig"
	providerhttpsms "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/sms/providerhttp"
	providertwilio "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/sms/providertwilio"
	config "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/smtp/config"
	featuressystem "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/system/features"
	actions "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/trigger/actions"
	accesstoken "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/user/accesstoken"
	grantuser "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/user/grant"
	humanuser "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/user/humanuser"
	keyuser "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/user/key"
	machineuser "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/user/machineuser"
	metadatauser "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/user/metadata"
	emailmessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/verify/emailmessagetext"
	emailotpmessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/verify/emailotpmessagetext"
	phonemessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/verify/phonemessagetext"
	smsotpmessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/verify/smsotpmessagetext"
	action "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/zitadel/action"
	domainzitadel "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/zitadel/domain"
	webkeyzitadel "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/namespaced/zitadel/webkey"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		executionevent.Setup,
		executionfunction.Setup,
		executionrequest.Setup,
		executionresponse.Setup,
		target.Setup,
		targetpublickey.Setup,
		webkey.Setup,
		api.Setup,
		key.Setup,
		oidc.Setup,
		saml.Setup,
		domainclaimedmessagetext.Setup,
		domainpolicy.Setup,
		initmessagetext.Setup,
		inviteusermessagetext.Setup,
		labelpolicy.Setup,
		lockoutpolicy.Setup,
		loginpolicy.Setup,
		logintexts.Setup,
		notificationpolicy.Setup,
		oidcsettings.Setup,
		passwordagepolicy.Setup,
		passwordchangemessagetext.Setup,
		passwordcomplexitypolicy.Setup,
		passwordlessregistrationmessagetext.Setup,
		passwordresetmessagetext.Setup,
		privacypolicy.Setup,
		securitysettings.Setup,
		verifyemailmessagetext.Setup,
		verifyemailotpmessagetext.Setup,
		verifyphonemessagetext.Setup,
		verifysmsotpmessagetext.Setup,
		claimedmessagetext.Setup,
		policy.Setup,
		providerhttp.Setup,
		providersmtp.Setup,
		apple.Setup,
		azuread.Setup,
		github.Setup,
		githubes.Setup,
		gitlab.Setup,
		gitlabselfhosted.Setup,
		google.Setup,
		ldap.Setup,
		oauth.Setup,
		oidcidp.Setup,
		samlidp.Setup,
		messagetext.Setup,
		customdomain.Setup,
		features.Setup,
		member.Setup,
		restrictions.Setup,
		secretgenerator.Setup,
		trusteddomain.Setup,
		policylabel.Setup,
		policylockout.Setup,
		policylogin.Setup,
		texts.Setup,
		policynotification.Setup,
		idpapple.Setup,
		idpazuread.Setup,
		idpgithub.Setup,
		idpgithubes.Setup,
		idpgitlab.Setup,
		idpgitlabselfhosted.Setup,
		idpgoogle.Setup,
		idpjwt.Setup,
		idpldap.Setup,
		idpoauth.Setup,
		idpoidc.Setup,
		idpsaml.Setup,
		memberorg.Setup,
		metadata.Setup,
		org.Setup,
		organization.Setup,
		domain.Setup,
		metadataorganization.Setup,
		agepolicy.Setup,
		changemessagetext.Setup,
		complexitypolicy.Setup,
		resetmessagetext.Setup,
		registrationmessagetext.Setup,
		policyprivacy.Setup,
		grant.Setup,
		grantmember.Setup,
		memberproject.Setup,
		project.Setup,
		role.Setup,
		providerconfig.Setup,
		providerhttpsms.Setup,
		providertwilio.Setup,
		config.Setup,
		featuressystem.Setup,
		actions.Setup,
		accesstoken.Setup,
		grantuser.Setup,
		humanuser.Setup,
		keyuser.Setup,
		machineuser.Setup,
		metadatauser.Setup,
		emailmessagetext.Setup,
		emailotpmessagetext.Setup,
		phonemessagetext.Setup,
		smsotpmessagetext.Setup,
		action.Setup,
		domainzitadel.Setup,
		webkeyzitadel.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		executionevent.SetupGated,
		executionfunction.SetupGated,
		executionrequest.SetupGated,
		executionresponse.SetupGated,
		target.SetupGated,
		targetpublickey.SetupGated,
		webkey.SetupGated,
		api.SetupGated,
		key.SetupGated,
		oidc.SetupGated,
		saml.SetupGated,
		domainclaimedmessagetext.SetupGated,
		domainpolicy.SetupGated,
		initmessagetext.SetupGated,
		inviteusermessagetext.SetupGated,
		labelpolicy.SetupGated,
		lockoutpolicy.SetupGated,
		loginpolicy.SetupGated,
		logintexts.SetupGated,
		notificationpolicy.SetupGated,
		oidcsettings.SetupGated,
		passwordagepolicy.SetupGated,
		passwordchangemessagetext.SetupGated,
		passwordcomplexitypolicy.SetupGated,
		passwordlessregistrationmessagetext.SetupGated,
		passwordresetmessagetext.SetupGated,
		privacypolicy.SetupGated,
		securitysettings.SetupGated,
		verifyemailmessagetext.SetupGated,
		verifyemailotpmessagetext.SetupGated,
		verifyphonemessagetext.SetupGated,
		verifysmsotpmessagetext.SetupGated,
		claimedmessagetext.SetupGated,
		policy.SetupGated,
		providerhttp.SetupGated,
		providersmtp.SetupGated,
		apple.SetupGated,
		azuread.SetupGated,
		github.SetupGated,
		githubes.SetupGated,
		gitlab.SetupGated,
		gitlabselfhosted.SetupGated,
		google.SetupGated,
		ldap.SetupGated,
		oauth.SetupGated,
		oidcidp.SetupGated,
		samlidp.SetupGated,
		messagetext.SetupGated,
		customdomain.SetupGated,
		features.SetupGated,
		member.SetupGated,
		restrictions.SetupGated,
		secretgenerator.SetupGated,
		trusteddomain.SetupGated,
		policylabel.SetupGated,
		policylockout.SetupGated,
		policylogin.SetupGated,
		texts.SetupGated,
		policynotification.SetupGated,
		idpapple.SetupGated,
		idpazuread.SetupGated,
		idpgithub.SetupGated,
		idpgithubes.SetupGated,
		idpgitlab.SetupGated,
		idpgitlabselfhosted.SetupGated,
		idpgoogle.SetupGated,
		idpjwt.SetupGated,
		idpldap.SetupGated,
		idpoauth.SetupGated,
		idpoidc.SetupGated,
		idpsaml.SetupGated,
		memberorg.SetupGated,
		metadata.SetupGated,
		org.SetupGated,
		organization.SetupGated,
		domain.SetupGated,
		metadataorganization.SetupGated,
		agepolicy.SetupGated,
		changemessagetext.SetupGated,
		complexitypolicy.SetupGated,
		resetmessagetext.SetupGated,
		registrationmessagetext.SetupGated,
		policyprivacy.SetupGated,
		grant.SetupGated,
		grantmember.SetupGated,
		memberproject.SetupGated,
		project.SetupGated,
		role.SetupGated,
		providerconfig.SetupGated,
		providerhttpsms.SetupGated,
		providertwilio.SetupGated,
		config.SetupGated,
		featuressystem.SetupGated,
		actions.SetupGated,
		accesstoken.SetupGated,
		grantuser.SetupGated,
		humanuser.SetupGated,
		keyuser.SetupGated,
		machineuser.SetupGated,
		metadatauser.SetupGated,
		emailmessagetext.SetupGated,
		emailotpmessagetext.SetupGated,
		phonemessagetext.SetupGated,
		smsotpmessagetext.SetupGated,
		action.SetupGated,
		domainzitadel.SetupGated,
		webkeyzitadel.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
