// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	executionevent "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/action/executionevent"
	executionfunction "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/action/executionfunction"
	executionrequest "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/action/executionrequest"
	executionresponse "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/action/executionresponse"
	target "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/action/target"
	targetpublickey "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/action/targetpublickey"
	webkey "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/active/webkey"
	api "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/application/api"
	key "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/application/key"
	oidc "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/application/oidc"
	saml "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/application/saml"
	domainclaimedmessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/default/domainclaimedmessagetext"
	domainpolicy "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/default/domainpolicy"
	initmessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/default/initmessagetext"
	inviteusermessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/default/inviteusermessagetext"
	labelpolicy "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/default/labelpolicy"
	lockoutpolicy "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/default/lockoutpolicy"
	loginpolicy "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/default/loginpolicy"
	logintexts "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/default/logintexts"
	notificationpolicy "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/default/notificationpolicy"
	oidcsettings "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/default/oidcsettings"
	passwordagepolicy "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/default/passwordagepolicy"
	passwordchangemessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/default/passwordchangemessagetext"
	passwordcomplexitypolicy "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/default/passwordcomplexitypolicy"
	passwordlessregistrationmessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/default/passwordlessregistrationmessagetext"
	passwordresetmessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/default/passwordresetmessagetext"
	privacypolicy "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/default/privacypolicy"
	securitysettings "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/default/securitysettings"
	verifyemailmessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/default/verifyemailmessagetext"
	verifyemailotpmessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/default/verifyemailotpmessagetext"
	verifyphonemessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/default/verifyphonemessagetext"
	verifysmsotpmessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/default/verifysmsotpmessagetext"
	claimedmessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/domain/claimedmessagetext"
	policy "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/domain/policy"
	providerhttp "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/email/providerhttp"
	providersmtp "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/email/providersmtp"
	apple "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/idp/apple"
	azuread "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/idp/azuread"
	github "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/idp/github"
	githubes "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/idp/githubes"
	gitlab "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/idp/gitlab"
	gitlabselfhosted "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/idp/gitlabselfhosted"
	google "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/idp/google"
	ldap "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/idp/ldap"
	oauth "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/idp/oauth"
	oidcidp "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/idp/oidc"
	samlidp "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/idp/saml"
	messagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/init/messagetext"
	customdomain "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/instance/customdomain"
	features "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/instance/features"
	member "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/instance/member"
	restrictions "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/instance/restrictions"
	secretgenerator "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/instance/secretgenerator"
	trusteddomain "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/instance/trusteddomain"
	policylabel "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/label/policy"
	policylockout "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/lockout/policy"
	policylogin "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/login/policy"
	texts "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/login/texts"
	policynotification "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/notification/policy"
	idpapple "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/org/idpapple"
	idpazuread "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/org/idpazuread"
	idpgithub "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/org/idpgithub"
	idpgithubes "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/org/idpgithubes"
	idpgitlab "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/org/idpgitlab"
	idpgitlabselfhosted "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/org/idpgitlabselfhosted"
	idpgoogle "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/org/idpgoogle"
	idpjwt "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/org/idpjwt"
	idpldap "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/org/idpldap"
	idpoauth "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/org/idpoauth"
	idpoidc "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/org/idpoidc"
	idpsaml "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/org/idpsaml"
	memberorg "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/org/member"
	metadata "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/org/metadata"
	org "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/org/org"
	organization "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/org/organization"
	domain "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/organization/domain"
	metadataorganization "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/organization/metadata"
	agepolicy "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/password/agepolicy"
	changemessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/password/changemessagetext"
	complexitypolicy "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/password/complexitypolicy"
	resetmessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/password/resetmessagetext"
	registrationmessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/passwordless/registrationmessagetext"
	policyprivacy "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/privacy/policy"
	grant "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/project/grant"
	grantmember "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/project/grantmember"
	memberproject "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/project/member"
	project "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/project/project"
	role "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/project/role"
	providerconfig "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/providerconfig"
	providerhttpsms "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/sms/providerhttp"
	providertwilio "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/sms/providertwilio"
	config "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/smtp/config"
	featuressystem "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/system/features"
	actions "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/trigger/actions"
	accesstoken "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/user/accesstoken"
	grantuser "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/user/grant"
	humanuser "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/user/humanuser"
	keyuser "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/user/key"
	machineuser "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/user/machineuser"
	metadatauser "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/user/metadata"
	emailmessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/verify/emailmessagetext"
	emailotpmessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/verify/emailotpmessagetext"
	phonemessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/verify/phonemessagetext"
	smsotpmessagetext "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/verify/smsotpmessagetext"
	action "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/zitadel/action"
	domainzitadel "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/zitadel/domain"
	webkeyzitadel "github.com/crossplane-contrib/provider-upjet-zitadel/internal/controller/cluster/zitadel/webkey"
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
