package weixinpay

//TODO: 关闭订单返回错误
/*
ORDERPAID	订单已支付	订单已支付，不能发起关单	订单已支付，不能发起关单，请当作已支付的正常交易
SYSTEMERROR	系统错误	系统错误	系统异常，请重新调用该API
ORDERNOTEXIST	订单不存在	订单系统不存在此订单	不需要关单，当作未提交的支付的订单
ORDERCLOSED	订单已关闭	订单已关闭，无法重复关闭	订单已关闭，无需继续调用
SIGNERROR	签名错误	参数签名结果不正确	请检查签名参数和方法是否都符合签名算法要求
REQUIRE_POST_METHOD	请使用post方法	未使用post传递参数 	请检查请求参数是否通过post方法提交
XML_FORMAT_ERROR	XML格式错误	XML格式错误	请检查XML参数格式是否正确
*/