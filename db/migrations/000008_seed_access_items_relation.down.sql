-- Remove seeded edges (adjust if you added other rows).
DELETE FROM ACCESS_ITEMS_RELATION WHERE PARENT_KEY IN (
    'cbe_to_cbe', 'mini_statement', 'bill_payment', 'voucher', 'chat', 'self_transfer',
    'scheduled_payment', 'merchant', 'utility', 'sitota', 'qr_payment', 'money_request',
    'split_bill', 'wallet', 'ips', 'map', 'news', 'exchange_rate', 'payment_verification',
    'donation', 'transaction_list', 'beneficiary', 'micro_finance', 'topup', 'mini_apps',
    'budget_management', 'micro_lending', 'locked_amount', 'event', 'quick_pay'
);
