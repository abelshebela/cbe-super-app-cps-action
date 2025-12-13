# CPS Actions Catalog

This document lists CPS Request Actions grouped by module and their derived event names for maker/checker phases.

Event name format:
- {module}.{action}.{phase}
- module and action are lowercase in events
- phase is one of: maker, checker

---

## Service
- CREATE_SERVICE_FEE → events: service.create_service_fee.maker / service.create_service_fee.checker
- UPDATE_SERVICE_FEE → events: service.update_service_fee.maker / service.update_service_fee.checker
- DELETE_SERVICE_FEE → events: service.delete_service_fee.maker / service.delete_service_fee.checker
- CREATE DAILY LIMIT → events: service.create daily limit.maker / service.create daily limit.checker
- UPDATE DAILY LIMIT → events: service.update daily limit.maker / service.update daily limit.checker
- DELETE DAILY LIMIT → events: service.delete daily limit.maker / service.delete daily limit.checker
- UPDATE_SERVICE_SINGLE_CAP → events: service.update_service_single_cap.maker / service.update_service_single_cap.checker
- UPDATE_SERVICE_TOTAL_CAP → events: service.update_service_total_cap.maker / service.update_service_total_cap.checker
- UPDATE_SERVICE_MIN_CAP → events: service.update_service_min_cap.maker / service.update_service_min_cap.checker

## ServicesCatalog
- CREATE_SERVICE → events: servicescatalog.create_service.maker / servicescatalog.create_service.checker
- UPDATE_SERVICE → events: servicescatalog.update_service.maker / servicescatalog.update_service.checker
- ENABLE_SERVICE → events: servicescatalog.enable_service.maker / servicescatalog.enable_service.checker
- DISABLE_SERVICE → events: servicescatalog.disable_service.maker / servicescatalog.disable_service.checker

## Account
- USER → events: account.user.maker / account.user.checker
- UPDATE_ACCOUNT_VALIDATION → events: account.update_account_validation.maker / account.update_account_validation.checker
- ENABLE_USER → events: account.enable_user.maker / account.enable_user.checker
- DISABLE_USER → events: account.disable_user.maker / account.disable_user.checker
- UPDATE_USER → events: account.update_user.maker / account.update_user.checker
- ARCHIVE_USER → events: account.archive_user.maker / account.archive_user.checker

## Event
- CREATE_EVENT → events: event.create_event.maker / event.create_event.checker
- DELETE_EVENT → events: event.delete_event.maker / event.delete_event.checker
- DISABLE_EVENT → events: event.disable_event.maker / event.disable_event.checker
- ENABLE_EVENT → events: event.enable_event.maker / event.enable_event.checker
- UPDATE_EVENT → events: event.update_event.maker / event.update_event.checker

## AmountBasedAuth
- CREATE_AMOUNT_BASED_AUTH → events: amountbasedauth.create_amount_based_auth.maker / amountbasedauth.create_amount_based_auth.checker
- UPDATE_AMOUNT_BASED_AUTH → events: amountbasedauth.update_amount_based_auth.maker / amountbasedauth.update_amount_based_auth.checker
- DELETE_AMOUNT_BASED_AUTH → events: amountbasedauth.delete_amount_based_auth.maker / amountbasedauth.delete_amount_based_auth.checker
- AUTH_TIER → events: amountbasedauth.auth_tier.maker / amountbasedauth.auth_tier.checker

## User
- USER → events: user.user.maker / user.user.checker
- ENABLE_USER → events: user.enable_user.maker / user.enable_user.checker
- DISABLE_USER → events: user.disable_user.maker / user.disable_user.checker
- UPDATE_USER → events: user.update_user.maker / user.update_user.checker
- ARCHIVE_USER → events: user.archive_user.maker / user.archive_user.checker

## BPSUser
- BPS_USER → events: bpsuser.bps_user.maker / bpsuser.bps_user.checker
- ENABLE_BPS_USER → events: bpsuser.enable_bps_user.maker / bpsuser.enable_bps_user.checker
- DISABLE_BPS_USER → events: bpsuser.disable_bps_user.maker / bpsuser.disable_bps_user.checker

## PermissionGroup
- PERMISSION_GROUP → events: permissiongroup.permission_group.maker / permissiongroup.permission_group.checker

## Department
- CREATE_DEPARTMENT → events: department.create_department.maker / department.create_department.checker
- UPDATE_DEPARTMENT → events: department.update_department.maker / department.update_department.checker
- ENABLE_DISABLE_DEPARTMENT → events: department.enable_disable_department.maker / department.enable_disable_department.checker

## ServiceFee
- CREATE_SERVICE_FEE → events: servicefee.create_service_fee.maker / servicefee.create_service_fee.checker
- UPDATE_SERVICE_FEE → events: servicefee.update_service_fee.maker / servicefee.update_service_fee.checker
- DELETE_SERVICE_FEE → events: servicefee.delete_service_fee.maker / servicefee.delete_service_fee.checker

## DailyLimit
- CREATE DAILY LIMIT → events: dailylimit.create daily limit.maker / dailylimit.create daily limit.checker
- UPDATE DAILY LIMIT → events: dailylimit.update daily limit.maker / dailylimit.update daily limit.checker
- DELETE DAILY LIMIT → events: dailylimit.delete daily limit.maker / dailylimit.delete daily limit.checker
- TOTAL_DAILY_LIMIT → events: dailylimit.total_daily_limit.maker / dailylimit.total_daily_limit.checker

## VAT
- UPDATE_VAT → events: vat.update_vat.maker / vat.update_vat.checker

## AuthTier
- AUTH_TIER → events: authtier.auth_tier.maker / authtier.auth_tier.checker

## Archive
- UPDATE_ARCHIVE_EXPIRY → events: archive.update_archive_expiry.maker / archive.update_archive_expiry.checker

## MinimumService
- UPDATE_MINIMUM_SERVICE → events: minimumservice.update_minimum_service.maker / minimumservice.update_minimum_service.checker

## ServiceRule
- UPDATE_SERVICE_RULE → events: servicerule.update_service_rule.maker / servicerule.update_service_rule.checker

## Total
- UPDATE_TOTAL → events: total.update_total.maker / total.update_total.checker

## AccessConfig
- UPDATE_ACCESS_CONFIG → events: accessconfig.update_access_config.maker / accessconfig.update_access_config.checker

## Branch
- ENABLE_SINGLE_BRANCH → events: branch.enable_single_branch.maker / branch.enable_single_branch.checker
- ENABLE_MULTI_USERS → events: branch.enable_multi_users.maker / branch.enable_multi_users.checker
- DISABLE_MULTI_USERS → events: branch.disable_multi_users.maker / branch.disable_multi_users.checker
- ENABLE_SINGLE_BRANCHES → events: branch.enable_single_branches.maker / branch.enable_single_branches.checker
- ENABLE_MULTI_BRANCHES → events: branch.enable_multi_branches.maker / branch.enable_multi_branches.checker

## Business
- CREATE_BUSINESS → events: business.create_business.maker / business.create_business.checker
- UPDATE_BUSINESS → events: business.update_business.maker / business.update_business.checker

## EventCategory
- CREATE_EVENT_CATEGORY → events: eventcategory.create_event_category.maker / eventcategory.create_event_category.checker
- UPDATE_EVENT_CATEGORY → events: eventcategory.update_event_category.maker / eventcategory.update_event_category.checker

## MiniAppMerchant
- CREATE_MINI_APP_MERCHANT → events: miniappmerchant.create_mini_app_merchant.maker / miniappmerchant.create_mini_app_merchant.checker
- UPDATE_MINI_APP_MERCHANT → events: miniappmerchant.update_mini_app_merchant.maker / miniappmerchant.update_mini_app_merchant.checker
- DELETE_MINI_APP_MERCHANT → events: miniappmerchant.delete_mini_app_merchant.maker / miniappmerchant.delete_mini_app_merchant.checker
- ENABLE_MINI_APP_MERCHANT → events: miniappmerchant.enable_mini_app_merchant.maker / miniappmerchant.enable_mini_app_merchant.checker
- DISABLE_MINI_APP_MERCHANT → events: miniappmerchant.disable_mini_app_merchant.maker / miniappmerchant.disable_mini_app_merchant.checker

## BlockTime
- UPDATE_BLOCK_TIME → events: blocktime.update_block_time.maker / blocktime.update_block_time.checker

## Password
- UPDATE_PASSWORD_RULE → events: password.update_password_rule.maker / password.update_password_rule.checker

## Permission
- CREATE_PERMISSION_GROUP → events: permission.create_permission_group.maker / permission.create_permission_group.checker
- DELETE_PERMISSION_GROUP → events: permission.delete_permission_group.maker / permission.delete_permission_group.checker
- UPDATE_PERMISSION_GROUP → events: permission.update_permission_group.maker / permission.update_permission_group.checker

## Avatar
- CREATE_AVATAR → events: avatar.create_avatar.maker / avatar.create_avatar.checker
- UPDATE_AVATAR → events: avatar.update_avatar.maker / avatar.update_avatar.checker
- ENABLE_AVATAR → events: avatar.enable_avatar.maker / avatar.enable_avatar.checker
- DISABLE_AVATAR → events: avatar.disable_avatar.maker / avatar.disable_avatar.checker
- DELETE_AVATAR → events: avatar.delete_avatar.maker / avatar.delete_avatar.checker

## Budget
- CREATE_BUDGET_COLOR → events: budget.create_budget_color.maker / budget.create_budget_color.checker
- UPDATE_BUDGET_COLOR → events: budget.update_budget_color.maker / budget.update_budget_color.checker
- DELETE_BUDGET_COLOR → events: budget.delete_budget_color.maker / budget.delete_budget_color.checker
- CREATE_BUDGET_ICON → events: budget.create_budget_icon.maker / budget.create_budget_icon.checker
- UPDATE_BUDGET_ICON → events: budget.update_budget_icon.maker / budget.update_budget_icon.checker
- DELETE_BUDGET_ICON → events: budget.delete_budget_icon.maker / budget.delete_budget_icon.checker

## Advert
- CREATE_ADVERT → events: advert.create_advert.maker / advert.create_advert.checker
- UPDATE_ADVERT → events: advert.update_advert.maker / advert.update_advert.checker
- ENABLE_ADVERT → events: advert.enable_advert.maker / advert.enable_advert.checker
- DISABLE_ADVERT → events: advert.disable_advert.maker / advert.disable_advert.checker
- DELETE_ADVERT → events: advert.delete_advert.maker / advert.delete_advert.checker

## Bank
- CREATE_BANK → events: bank.create_bank.maker / bank.create_bank.checker
- UPDATE_BANK → events: bank.update_bank.maker / bank.update_bank.checker
- DELETE_BANK → events: bank.delete_bank.maker / bank.delete_bank.checker
- ENABLE_DISABLE_BANK → events: bank.enable_disable_bank.maker / bank.enable_disable_bank.checker
- UPDATE_BANK_LOGO → events: bank.update_bank_logo.maker / bank.update_bank_logo.checker
- ENABLE_BANK → events: bank.enable_bank.maker / bank.enable_bank.checker
- DISABLE_BANK → events: bank.disable_bank.maker / bank.disable_bank.checker

## Wallet
- CREATE_WALLET → events: wallet.create_wallet.maker / wallet.create_wallet.checker
- UPDATE_WALLET → events: wallet.update_wallet.maker / wallet.update_wallet.checker
- DELETE_WALLET → events: wallet.delete_wallet.maker / wallet.delete_wallet.checker
- ENABLE_WALLET → events: wallet.enable_wallet.maker / wallet.enable_wallet.checker
- DISABLE_WALLET → events: wallet.disable_wallet.maker / wallet.disable_wallet.checker

## Topup
- CREATE_TOPUP → events: topup.create_topup.maker / topup.create_topup.checker
- UPDATE_TOPUP → events: topup.update_topup.maker / topup.update_topup.checker
- DELETE_TOPUP → events: topup.delete_topup.maker / topup.delete_topup.checker
- ENABLE_TOPUP → events: topup.enable_topup.maker / topup.enable_topup.checker
- DISABLE_TOPUP → events: topup.disable_topup.maker / topup.disable_topup.checker

## Validation
- CREATE_VALIDATION → events: validation.create_validation.maker / validation.create_validation.checker
- UPDATE_VALIDATION → events: validation.update_validation.maker / validation.update_validation.checker
- DELETE_VALIDATION → events: validation.delete_validation.maker / validation.delete_validation.checker

## Block
- BLOCK_USER → events: block.block_user.maker / block.block_user.checker
- DISABLE_SINGLE_BRANCH → events: block.disable_single_branch.maker / block.disable_single_branch.checker
- ENABLE_SINGLE_BRANCH → events: block.enable_single_branch.maker / block.enable_single_branch.checker
- DISABLE_MULTI_BRANCHES → events: block.disable_multi_branches.maker / block.disable_multi_branches.checker
- ENABLE_MULTI_BRANCHES → events: block.enable_multi_branches.maker / block.enable_multi_branches.checker
- ENABLE_BRANCHES → events: block.enable_branches.maker / block.enable_branches.checker
- DISABLE_BRANCHES → events: block.disable_branches.maker / block.disable_branches.checker
- ENABLE_REGIONS → events: block.enable_regions.maker / block.enable_regions.checker
- DISABLE_REGIONS → events: block.disable_regions.maker / block.disable_regions.checker
- ENABLE_DISTRICTS → events: block.enable_districts.maker / block.enable_districts.checker
- DISABLE_DISTRICTS → events: block.disable_districts.maker / block.disable_districts.checker
- ENABLE_CITIES → events: block.enable_cities.maker / block.enable_cities.checker
- DISABLE_CITIES → events: block.disable_cities.maker / block.disable_cities.checker

## BudgetCategory
- CREATE_BUDGET_CATEGORY → events: budgetcategory.create_budget_category.maker / budgetcategory.create_budget_category.checker
- UPDATE_BUDGET_CATEGORY → events: budgetcategory.update_budget_category.maker / budgetcategory.update_budget_category.checker
- DELETE_BUDGET_CATEGORY → events: budgetcategory.delete_budget_category.maker / budgetcategory.delete_budget_category.checker
- DISABLE_BUDGET_CATEGORY → events: budgetcategory.disable_budget_category.maker / budgetcategory.disable_budget_category.checker
- ENABLE_BUDGET_CATEGORY → events: budgetcategory.enable_budget_category.maker / budgetcategory.enable_budget_category.checker

## UnlinkDevice
- UNLINK_DEVICE → events: unlinkdevice.unlink_device.maker / unlinkdevice.unlink_device.checker
- UNLINK_USER → events: unlinkdevice.unlink_user.maker / unlinkdevice.unlink_user.checker

## HQ
- UPDATE_HQ_BLOCK_TIME → events: hq.update_hq_block_time.maker / hq.update_hq_block_time.checker
- UPDATE_HQ_ARCHIVE_TIME → events: hq.update_hq_archive_time.maker / hq.update_hq_archive_time.checker
- UPDATE_PASSWORD_EXPIRY → events: hq.update_password_expiry.maker / hq.update_password_expiry.checker

## Fayda
- DISABLE_FAYDA_ACCOUNT → events: fayda.disable_fayda_account.maker / fayda.disable_fayda_account.checker
- ENABLE_FAYDA_ACCOUNT → events: fayda.enable_fayda_account.maker / fayda.enable_fayda_account.checker

## CPSUser
- CPS_USER_CREATE → events: cpsuser.cps_user_create.maker / cpsuser.cps_user_create.checker
- CPS_USER_UPDATE → events: cpsuser.cps_user_update.maker / cpsuser.cps_user_update.checker
- CPS_USER_DELETE → events: cpsuser.cps_user_delete.maker / cpsuser.cps_user_delete.checker
- CPS_USER_ENABLE → events: cpsuser.cps_user_enable.maker / cpsuser.cps_user_enable.checker
- CPS_USER_DISABLE → events: cpsuser.cps_user_disable.maker / cpsuser.cps_user_disable.checker

## BulkService
- BULK_SERVICE_ENABLE → events: bulkservice.bulk_service_enable.maker / bulkservice.bulk_service_enable.checker
- BULK_SERVICE_DISABLE → events: bulkservice.bulk_service_disable.maker / bulkservice.bulk_service_disable.checker

## MiniApp
- CREATE_MINI_APP → events: miniapp.create_mini_app.maker / miniapp.create_mini_app.checker
- UPDATE_MINI_APP → events: miniapp.update_mini_app.maker / miniapp.update_mini_app.checker
- DELETE_MINI_APP → events: miniapp.delete_mini_app.maker / miniapp.delete_mini_app.checker
- ENABLE_MINI_APP → events: miniapp.enable_mini_app.maker / miniapp.enable_mini_app.checker
- DISABLE_MINI_APP → events: miniapp.disable_mini_app.maker / miniapp.disable_mini_app.checker

## Notification
- CREATE_PUBLIC_NOTIFICATION → events: notification.create_public_notification.maker / notification.create_public_notification.checker
- UPDATE_PUBLIC_NOTIFICATION → events: notification.update_public_notification.maker / notification.update_public_notification.checker
- DELETE_NOTIFICATION → events: notification.delete_notification.maker / notification.delete_notification.checker
- ENABLE_NOTIFICATION → events: notification.enable_notification.maker / notification.enable_notification.checker
- DISABLE_NOTIFICATION → events: notification.disable_notification.maker / notification.disable_notification.checker
- MARK_NOTIFICATION_AS_SEEN → events: notification.mark_notification_as_seen.maker / notification.mark_notification_as_seen.checker

## BankVault
- CREATE_BANK_VAULT → events: bankvault.create_bank_vault.maker / bankvault.create_bank_vault.checker
- UPDATE_BANK_VAULT → events: bankvault.update_bank_vault.maker / bankvault.update_bank_vault.checker
- DELETE_BANK_VAULT → events: bankvault.delete_bank_vault.maker / bankvault.delete_bank_vault.checker
- ENABLE_BANK_VAULT → events: bankvault.enable_bank_vault.maker / bankvault.enable_bank_vault.checker
- DISABLE_BANK_VAULT → events: bankvault.disable_bank_vault.maker / bankvault.disable_bank_vault.checker

## ProductCode
- UPDATE_PRODUCT_CODE → events: productcode.update_product_code.maker / productcode.update_product_code.checker

## Donation
- CREATE_DONATION → events: donation.create_donation.maker / donation.create_donation.checker
- UPDATE_DONATION → events: donation.update_donation.maker / donation.update_donation.checker
- DISABLE_DONATION → events: donation.disable_donation.maker / donation.disable_donation.checker
- ADD_DONATION_IMAGE → events: donation.add_donation_image.maker / donation.add_donation_image.checker
- UPDATE_DONATION_IMAGE → events: donation.update_donation_image.maker / donation.update_donation_image.checker
- DELETE_DONATION_IMAGE → events: donation.delete_donation_image.maker / donation.delete_donation_image.checker
- ENABLE_DONATION → events: donation.enable_donation.maker / donation.enable_donation.checker

## donationCategory
- CREATE_DONATION_CATEGORY → events: donationcategory.create_donation_category.maker / donationcategory.create_donation_category.checker
- UPDATE_DONATION_CATEGORY → events: donationcategory.update_donation_category.maker / donationcategory.update_donation_category.checker
- DISABLE_DONATION_CATEGORY → events: donationcategory.disable_donation_category.maker / donationcategory.disable_donation_category.checker
- ENABLE_DONATION_CATEGORY → events: donationcategory.enable_donation_category.maker / donationcategory.enable_donation_category.checker

## donationCompany
- CREATE_DONATION_COMPANY → events: donationcompany.create_donation_company.maker / donationcompany.create_donation_company.checker
- UPDATE_DONATION_COMPANY → events: donationcompany.update_donation_company.maker / donationcompany.update_donation_company.checker
- ENABLE_DONATION_COMPANY → events: donationcompany.enable_donation_company.maker / donationcompany.enable_donation_company.checker
- DISABLE_DONATION_COMPANY → events: donationcompany.disable_donation_company.maker / donationcompany.disable_donation_company.checker

## VaultGroupCategory
- CREATE_VAULT_GROUP_CATEGORY → events: vaultgroupcategory.create_vault_group_category.maker / vaultgroupcategory.create_vault_group_category.checker
- UPDATE_VAULT_GROUP_CATEGORY → events: vaultgroupcategory.update_vault_group_category.maker / vaultgroupcategory.update_vault_group_category.checker
- DELETE_VAULT_GROUP_CATEGORY → events: vaultgroupcategory.delete_vault_group_category.maker / vaultgroupcategory.delete_vault_group_category.checker
- ENABLE_VAULT_GROUP_CATEGORY → events: vaultgroupcategory.enable_vault_group_category.maker / vaultgroupcategory.enable_vault_group_category.checker
- DISABLE_VAULT_GROUP_CATEGORY → events: vaultgroupcategory.disable_vault_group_category.maker / vaultgroupcategory.disable_vault_group_category.checker

## KYCVerifier
- UPDATE_KYC_VERIFIER → events: kycverifier.update_kyc_verifier.maker / kycverifier.update_kyc_verifier.checker
- APPROVE_KYC → events: kycverifier.approve_kyc.maker / kycverifier.approve_kyc.checker

## article
- CREATE_ARTICLE → events: article.create_article.maker / article.create_article.checker
- UPDATE_ARTICLE → events: article.update_article.maker / article.update_article.checker
- DELETE_ARTICLE → events: article.delete_article.maker / article.delete_article.checker
- ENABLE_ARTICLE → events: article.enable_article.maker / article.enable_article.checker
- DISABLE_ARTICLE → events: article.disable_article.maker / article.disable_article.checker

## articleCategory
- CREATE_ARTICLE_CATEGORY → events: articlecategory.create_article_category.maker / articlecategory.create_article_category.checker
- UPDATE_ARTICLE_CATEGORY → events: articlecategory.update_article_category.maker / articlecategory.update_article_category.checker
- DELETE_ARTICLE_CATEGORY → events: articlecategory.delete_article_category.maker / articlecategory.delete_article_category.checker
- ENABLE_ARTICLE_CATEGORY → events: articlecategory.enable_article_category.maker / articlecategory.enable_article_category.checker
- DISABLE_ARTICLE_CATEGORY → events: articlecategory.disable_article_category.maker / articlecategory.disable_article_category.checker

## short_video
- CREATE_SHORT_VIDEO → events: short_video.create_short_video.maker / short_video.create_short_video.checker
- UPDATE_SHORT_VIDEO → events: short_video.update_short_video.maker / short_video.update_short_video.checker
- DELETE_SHORT_VIDEO → events: short_video.delete_short_video.maker / short_video.delete_short_video.checker
- ENABLE_SHORT_VIDEO → events: short_video.enable_short_video.maker / short_video.enable_short_video.checker
- DISABLE_SHORT_VIDEO → events: short_video.disable_short_video.maker / short_video.disable_short_video.checker

## customer
- ENABLE_DISABLE_CUSTOMER → events: customer.enable_disable_customer.maker / customer.enable_disable_customer.checker
- APPROVE_FAYDA_CUSTOMER → events: customer.approve_fayda_customer.maker / customer.approve_fayda_customer.checker

## news_category
- CREATE_NEWS_CATEGORY → events: news_category.create_news_category.maker / news_category.create_news_category.checker
- UPDATE_NEWS_CATEGORY → events: news_category.update_news_category.maker / news_category.update_news_category.checker
- DELETE_NEWS_CATEGORY → events: news_category.delete_news_category.maker / news_category.delete_news_category.checker

## news_tag
- CREATE_NEWS_TAG → events: news_tag.create_news_tag.maker / news_tag.create_news_tag.checker
- UPDATE_NEWS_TAG → events: news_tag.update_news_tag.maker / news_tag.update_news_tag.checker
- ENABLE_NEWS_TAG → events: news_tag.enable_news_tag.maker / news_tag.enable_news_tag.checker
- DISABLE_NEWS_TAG → events: news_tag.disable_news_tag.maker / news_tag.disable_news_tag.checker
- DELETE_NEWS_TAG → events: news_tag.delete_news_tag.maker / news_tag.delete_news_tag.checker

## ActionRole
- CREATE_ACTION_ROLE → events: actionrole.create_action_role.maker / actionrole.create_action_role.checker
- UPDATE_ACTION_ROLE → events: actionrole.update_action_role.maker / actionrole.update_action_role.checker
- ENABLE_ACTION_ROLE → events: actionrole.enable_action_role.maker / actionrole.enable_action_role.checker
- DISABLE_ACTION_ROLE → events: actionrole.disable_action_role.maker / actionrole.disable_action_role.checker

## CpsActionRole
- CREATE_CPS_ACTION_ROLE → events: cpsactionrole.create_cps_action_role.maker / cpsactionrole.create_cps_action_role.checker
- UPDATE_CPS_ACTION_ROLE → events: cpsactionrole.update_cps_action_role.maker / cpsactionrole.update_cps_action_role.checker
- ENABLE_CPS_ACTION_ROLE → events: cpsactionrole.enable_cps_action_role.maker / cpsactionrole.enable_cps_action_role.checker
- DISABLE_CPS_ACTION_ROLE → events: cpsactionrole.disable_cps_action_role.maker / cpsactionrole.disable_cps_action_role.checker

## DeviceVersion
- CREATE_DEVICE_VERSION → events: deviceversion.create_device_version.maker / deviceversion.create_device_version.checker
- UPDATE_DEVICE_VERSION → events: deviceversion.update_device_version.maker / deviceversion.update_device_version.checker
- ENABLE_DEVICE_VERSION → events: deviceversion.enable_device_version.maker / deviceversion.enable_device_version.checker
- DISABLE_DEVICE_VERSION → events: deviceversion.disable_device_version.maker / deviceversion.disable_device_version.checker
- DELETE_DEVICE_VERSION → events: deviceversion.delete_device_version.maker / deviceversion.delete_device_version.checker
- ENABLE_DISABLE_DEVICE_VERSION → events: deviceversion.enable_disable_device_version.maker / deviceversion.enable_disable_device_version.checker

## MiniAppCategory
- CREATE_MINI_APP_CATEGORY → events: miniappcategory.create_mini_app_category.maker / miniappcategory.create_mini_app_category.checker
- UPDATE_MINI_APP_CATEGORY → events: miniappcategory.update_mini_app_category.maker / miniappcategory.update_mini_app_category.checker
- DELETE_MINI_APP_CATEGORY → events: miniappcategory.delete_mini_app_category.maker / miniappcategory.delete_mini_app_category.checker
- ENABLE_MINI_APP_CATEGORY → events: miniappcategory.enable_mini_app_category.maker / miniappcategory.enable_mini_app_category.checker
- DISABLE_MINI_APP_CATEGORY → events: miniappcategory.disable_mini_app_category.maker / miniappcategory.disable_mini_app_category.checker

---

Notes:
- This file is maintained from RequestActionGroups in code. If you add a new action/group, update this doc accordingly.
- Event names are a convention to tag Kafka logs or metrics per module/action/phase.
