package module

import (
	"fmt"
	http "github.com/ProjectAthenaa/sonic-core/fasttls"
	"github.com/ProjectAthenaa/sonic-core/fasttls/cookiejar"
	"github.com/ProjectAthenaa/sonic-core/protos/module"
)

func (tk *Task) ATC() {
	req, err := tk.NewRequest("POST", "https://www.newbalance.com/on/demandware.store/Sites-NBUS-Site/en_US/Cart-AddProduct", []byte(`pid=194768477287&quantity=1&estimatedDelivery=Estimated+delivery+-+2-5+business+days+once+shipped&options=%5B%5D`))
	//req, err := tk.NewRequest("POST", "https://www.newbalance.com/on/demandware.store/Sites-NBUS-Site/en_US/Cart-AddProduct", []byte(fmt.Sprintf("pid=%s&quantity=1&estimatedDelivery=Estimated+delivery+-+2-5+business+days+once+shipped&options=%%5B%%5D", tk.VariantId)))
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "could not create atc req")
		tk.Stop()
		return
	}
	req.Headers = tk.GenerateDefaultHeaders(fmt.Sprintf(`https://www.newbalance.com/pd/~/%s.html`, tk.PID))

	cookiejar.ReleaseCookieJar(tk.FastClient.Jar)
	tk.FastClient.Jar = nil

	//req.Headers[`x-dtpc`] = []string{`2$500472369_531h28vPUNGQCPHUHIIUBHMKQIJELBUOKKUAIFM-0e2`}
	//req.Headers["Accept"] = []string{"*/*"}

	shapeheaders, err := tk.getShapeHeaders()
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "could not retrieve shape headers")
		tk.Stop()
		return
	}
	for k, v := range shapeheaders {
		req.Headers[k] = []string{v}
	}

	req.Headers = http.Headers{
		`User-Agent`: []string{`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36`},
		`Host`: []string{`www.newbalance.com`},
		`Content-Type`: []string{`application/x-www-form-urlencoded; charset=UTF-8`},
		`Content-Length`: []string{`61`},
		`Accept-Encoding`: []string{`gzip, deflate, br`},
		`X-Requested-With`: []string{`XMLHttpRequest`},
		`Cache-Control`: []string{`no-cache`},
		`PRMeACzwub-f`: []string{`A1sVkK17AQAAhuyai-lk7MfBAWb331YtpD_FlT5OE6r3qEPcCW3hzcGTw3ANAUha_m4KT-IXwH8AAAAAAAAAAA==`},
		`PRMeACzwub-z`: []string{`q`},
		`PRMeACzwub-a`: []string{`x1d23s2ZJyWxml22e-dzNMn8I1i2tYb=8bb2BuIVdYF3sALZZl8BqCkkomnwRvlpSHJY=W0G5CRa1vgs2zbAmslmGVw1TWZEXeV0A-jLwWldy2nSuk8ilKQ0eEsFy5ILgu2MurbntKb4Ehw8TGusWCOnUsTMHTNXnR4qWAQ8pCZVUDnwwzWvEa5juh7kvtwXxm4TySeLJUcNIx_jp2q10RIWyUpzIpfmXWDhKjM=NUlmwPvIofbdJs1OsZomID_c5viTtTJNug3ajoYxIJJCXOLWWwoJpxNpG2mxwF=HX23O_3WcAEgKY_5mb080I-cDGRV5nTN4GZNEi99Xh5s93saEVUYEJeaqxGcgGueDOjy_HQGLMgY0hG7B5AQftrNngSGRWqNmyFtoBfCMMIzUtE_hVwIr3ob9YXzfSS3iBvR0YhLJwvyiQg=_frG7QaaJq4iggaWumkYYGDYK0qeFxfHiU8oM2roQNht8urY9=ESkaZ7fF3OmRADJl41EAbnFGVw3RPJLsDiDh2b=e=leh4__E7YHDjEykA2GS4nF0MSOOOb598lvV4qBqyxdbu1HYmU5IOnFyWiLzR7H07OMUaBJx3ouMRkZXjweo9gOf5oZBqHLnrsPtyrO7L73=MHl33e_cT41j=JF1Xk3BFXELb8VEZiND1WiismvheQukULt1yIk--=Ij_UlhkublWkqyi_J4CJS7OXZpFzJhKbZWj7ZSK1=_Hpacf1OIPzOl33tPkuWBk0F2NUqoS9B5UI42OT3rjII4-2XFWJpyaZNVxRyp5ptOGEHEGysF1Rg=Cdl0VYj=wyjzc5YyqZLeSIz9bja7yLx7URcD2pAO2pU8vQiKqFzifegAUoreXexT3YAjRBsOROXW79lkBVlGj8qMZ2TVBqx4tM3iPqYi04y91JfZRTvfPSQkr3skkwRGVyfxvJ4jWb3ksFqUnCDM=xAdzOUwA3eJfOCUAz__Oh2wvGSdKQf-Ps9S5=pwXd2jSsMStACPKurDQaoXt1mIFuWIwyHfw3vUN-AnCbKjCujuRIXFnegrw5ow2RkkJ_jgLq7CxUKf3iVzW_A7qivKrqa_dMUl40N1=91M1kHU8s0jPtT4XzW4JQ5ma4IlskxD3AKhnH4FA3ks9ZHSAJpOeomA1XyHS1UrIRKd25Pkowiz8UO3JlcByVau4fTSH8VqSChDuVJDIiXj-ktrZSrCUasLaUBKgHF8E=p7cvbNOXugJT73ygWiWgkHtACIPj-VxzwOjdZhYeR2OKyeM07RfnqnxxvMGqMuv-GWrkD4z1np5bnp-bYyMvNDJDNleiyUP1C13eKv8BWamMNJw-jbHIW5ubnvtoRpT9yCaO2ZIMECUJ19IFdoaKcdnujbCp3OdXTUS2-hxOhCx52RsU70cds-WTqB4uEBvNIzf8DM-du4bMmPxMjpmp-MjpgBC3AbhgFHX2bYTxMQmrUoIhUzSq9VRB_VuHEhvW0i5NNMugLjg7gFQo4ufK0Zejm7=OD4P1tTwusgiZ2TpqtdKatXC_50iXJIeidv=BMkMXEtuAWlbMJj0SN8Acylr8CKF7OZbtYCCGb4AiMYt_ziyJ2uM270SPipUt1U0adhOwzn9ho7zBkgO4-Ks3EwMeTQ5r55W=wrEUYgR=Hm9UjGmRs-YX0EjknAvvGfQD1BeiQLC4Me=QPNApiO2YikUFqR5sEtpvQVVmjVnBGebQ80e4FNTRCA-Y9MJKF1crE=3KtfyAzmghTi30Jel3g_-ov7hMLZ97BxsADSjD7fxp5bLEELijMFzA5f_8dhrmGuKkU_aVP_K3tQfGxUNvyOjDY7u3GyBbIFd5cdZvL3gqkhgwo55xBbhFKANGxSmi4RGEAKzvhJQ7Zj7s_pJaB5nXQ8UYLy4-niN0cppmV2SJZXv8RflopwGB18mNJZsLjC3YuSDt5tJeVEYW9znH_LKkJkYFVwmnb-zbUyxiIPx353X14vctpSboLWyShDl1W90dS5=UJ7hnHITPhLEAN3zafT50iMzFinHZlB=amZKpNwrMBmO1w373rPzQcYgHdlbr3gK=f9SPq74sQuZYtxDq-O1b1NIntO1aVUMDFyCBKqobvMWlG5krGQtuVCumrX4nnokniQzdghAwWx-GnN=BBz7prVBCctjJYWPW291YSdAUPnx8QdcqsoeY0xWkLklpK7IDNZ_zN9EI7ExMbGizB5xAD=lGbGnUO01I3l_h0I3Wy-ck_m1SzdbuKLqVKqZB-L0LufveQCwaemxpXbVjPdEmYAgV00GkSo0IsU1IAo1l7u9m0xl9l1kO_VB5PeRuX8TtL2mgfw3sZx5f2=kyTO_lLTVzfApIIeIj5Y=THND2yMvR1GGatdmKCspIxjkLWNwqwNQGkrGz_tA4TXzfzdyEXcYyk9HUWapCptAd=7k=cWP0uBrIqjvpZwag8ECNi37Dnfqq7BWXoMsE7onEEBV=K5iYEOCvR4A7RdfwyZda5aBqhPLhgtRQ=KMjwurQHxAvOWbDRFlsdqgXEbTpbc1WyQ47SAzR7MyPbyIizeo00QdR48AkH9v-G38LlMnnOlZQLroGLrb17cDFckoFtAkqI-UBfl-1P3u11-Zog8mFXwmn_vZAfmrMHvYkfSeUodGszlHjWRxbKUQvgkMOElOUsbd9pWGq1hVvth3QBJP_pwsRT9z4f_F0W2pISs2hCnvgSuGUTTl8M3vnGFJY_=B2LPt-QfIEf2q9vY4VzVuOAPx3YvmiTQbizIiT1vbkHDwUPw2OR1Qy30qD5y5y4EYEqanGn8oHrNZ2Z9micq3xCse3oN9ekVR3Gwi04ng-Jo5qvatuqRPnfN92MiXFGwStLTopdqORkxVF1ILOV1G8r0hOhQ0cfiQ2t8hTWnIrZU85XsKxsKKAmCmda1SUVe8BX8bvNcaKwZzm4BxtEn2mZmfUqAj-5rgJz95YFrb5MPrwABh9P3vAxcwKWtGnlZy==2jqS_tlWwaniOKkpL-npeuwM8RS3ZbhWtMGYfJHDCwzJhyGoGNXcZ3wepClqkksowDHZPCsAv2ETT-jhBrgaMhHq0liDLScEScZWetYzQXaaIjO8_3e4ATHxj=Z89-huOxxaL-WGZjgV2AkJA_gVUNW-N-e2YHKumZrZMs7w1dBGjb3vg_TLl5sGaWO=K4dyrrxxBeUfDvvCtHykMcgsg=IFdvt9ssEY=oQ7fE3ARs=TezIQawK2jpgHCyk0VGiZpX2riOrJvKGQyTBNvTc5qiORQ2TtrKN7Jesinq9l8UID9hL=Qt-8zrN=LMCy1UmKeOrX7o4IGbganwXFvUkIT8YVSgi-9GlHrds5OSII_QxrV=JTPUT4x-bIwyedQj-LKmyFwD_Oe5FfQ8mm0IlkHm=YC8v=9x1EaD11vBYgAaH9VRIdI0wzvpvGjPV8oqYQxtrYhEEyzQ03sm3Or85zv1khDMcjYZOe9AbDVHG8x4SC15bDhiqmLpAUpZYItfWKtkjGG3nbBh3xgeUDGR2dXJh3CGB351LTXc3k=vfh9i5NeUzlpaFZ1dcDnDGu7n0Y=sB-ENLMPjvw2-F4INGUc_vbdcwv=WrE9FFCYIui0=XDLO=crD-DaQSDpdsSUX9n_RJBEmvyYhc_uDlRbbuSgMqzY=tdnKtsNx9ANtKfcpea2GT-X3XigDpr7CQtJBlYy7wKtykMPe41z=zMAxa00LSLuzVMsHMTgHV9NW29R31-7LtdFpdbhfNKFxC=u-_EzA3W98hTY9GdkCVSeGmeN84zGXo2eP=kIzjR2yExqgeLqMVITfW7PNLvkK4nrOueO1r4-Nir_Ydd5xshx2ah2gF7ByKdCPC59YhQvEQoG81NNKbexfVhAF37gh5t4wqKHDn=kDBv05PMWOkJ9M1i1Liinb4hHGPux2Cas-3Y3q9YFDQv4BDFKKdrcKm2JkO51Hq1IHWlpTSe5TP0nLFz915cRGh2_g7uwUn0GU2XHj1pW7XVRFUYUGwZGT_KybX7PGr4FLbDyFLcHRvlxf0crHC2nPMgy4VWxqdRM7VSqW3Y4Uonfd0YoU5LVSs5WUBpeMq0m9oJamnlw1Ko07CPUY848qPw9BkTBl=k=u=xwdeqZsmphdkFkp3iSk_fAerGGWgc=vfvFuDsSc3jqc_pS9FT=57_DyOMOq3TUNqzv8QCbu-B5lYiursdTgO1hk000Bg-hlwG-HtNEMMu0zBDcv_sevnU8CO99yV8MR0x=FK--FNbppnSjAE1fGE-u8ygo8eEhbsGdyv=Eb3UUfTu_lS3ab7CBVPUNIj2VLCHTP3-zKx-9qR9UMqf1TSmqEIWXaZzFGQmGMDjLe2LgdzBuuwNx-7qHcLhEP0yJX1XDv2xE5TyC80pQSx_BBgEmUCDVw3HsBJKXAmeOt=lY=vwFUFP-7JO3YlevIF_AB0UAzK910jWCDIlAK-oUTUC2ClEwGinReUVwOtX-FRPq3DeTZQU1U40jNXdk=a4uZl0B_Q5kT32aOMt1qSgoA9Uy4fBvvY-iXdeh8BKYBXTOc2zL2=uAI0QH1zbwsm=eflXHUkZrQIoXyOlvIYUEkuvIgvy5s3v2wSmBQ9vu0V7EXFri5llfi_dDg=acA0m0XivoGKO=hnlgttpg4gsDkqNspPk4WnkMqBxU-cRu10EoDnokB-ESeEVUoY=oD-dKchWVuaHYf9rO1DWpODPrjPUOEe_LwrHNUQUYk2KvoocjEeihJ9mtq5=sKUSZL4VUMI0kRdLWKnLx1VTqmFZjnW1eOs=54B0ldgfWmr-WEVknCpIZs7QY=s2eTut71LR7FlvLr4dcO_iy3y=b_d5IjvnEkCh_30f`},
		`Accept-Language`: []string{`en-us`},
		`Sec-Ch-Ua-Mobile`: []string{`?0`},
		`Origin`: []string{`https://www.newbalance.com`},
		`PRMeACzwub-c`: []string{`AIAMdq17AQAACxoF274Ur-riZFOheETc76aY6fnCqjPbTiqKlTIGd9Znv_0U`},
		`X-Dtpc`: []string{`2$500472369_531h28vPUNGQCPHUHIIUBHMKQIJELBUOKKUAIFM-0e2`},
		`PRMeACzwub-d`: []string{`ABaChIjBDKGNgUGAQZIQhISi0eIApJmBDgAyBnfWZ7_9FAAAAABy589kAAm-yEEouwIc7BnwBQ1ivg0`},
		`PRMeACzwub-b`: []string{`74iogm`},
		`Accept`: []string{`application/json`},
		`Sec-Fetch-Site`: []string{`same-site`},
		`Sec-Fetch-Dest`: []string{`empty`},
		`Sec-Fetch-Mode`: []string{`cors`},
		`Connection`: []string{`keep-alive`},
		`Sec-Ch-Ua`: []string{`"Chromium";v="91", " Not A;Brand";v="99", "Google Chrome";v="91"`},
		`Referer`: []string{`https://www.newbalance.com/pd/~/MCRZRV1-35865.html`},
		`Pragma`: []string{`no-cache`},
		"Cookie": []string{`dwsid=8YorJskuxsnSpTXfas-Wcp_3fesfofMxHtYA1fhCRVX-QcKMwt6Y0mE3Ogo8lLg9rNyNY2PkE8dimROk_sX-Ww==; AKA_A2=A; dwac_4eef82b2a11317de54af2b5132=1zpRWj8FoyUNw1Sp4tiZO7lHB5fpipbdwbA%3D|dw-only|||USD|false|US%2FEastern|true; cqcid=abK8spjqw513eTaXjAPoQRZXAo; cquid=||; sid=1zpRWj8FoyUNw1Sp4tiZO7lHB5fpipbdwbA; dwanonymous_b46d190781ef77bb66faac87f06d52c0=abK8spjqw513eTaXjAPoQRZXAo; __cq_dnt=0; dw_dnt=0; optimizelyEndUserId=0efa3b17cd640000b0833261e3010000a0492200; _cs_c=1; __cq_bc=%7B%22aagi-NBUS%22%3A%5B%7B%22id%22%3A%22MCRZRV1-35865%22%7D%5D%7D; __cq_uuid=abK8spjqw513eTaXjAPoQRZXAo; _adobe_cp_scopen=open; utag_vnum=1633292465410&vn=1; s_fid=52EB866B2AFC1BF4-07DBA160AD2E2976; s_cc=true; s_vi=[CS]v1|309941D89DB69BA2-40001B74405C6059[CE]; _scid=4a14107d-8556-4ca1-b2c6-d273b6482d38; __pdst=2005412fecc34a2e9d47bfc1d404514f; _cavisit=17bad526f29|; _caid=89a098b4-b6a4-4d94-b91f-dbb6f5b1004a; _ga=GA1.2.1775671219.1630700466; _gid=GA1.2.938499465.1630700466; scarab.mayAdd=%5B%7B%22i%22%3A%22MCRZRV1-35865%22%7D%5D; rxVisitor=16307004659907GAV627UN0SU7LG34O5QKUDIBVB77QVT; dtSa=-; _pin_unauth=dWlkPU5EUXdaR1k0TkRjdE5XRm1OeTAwT0RjeUxUaGxPRGt0T1dObU9XWXpNbUZrTVdWaA; _gcl_au=1.1.246205982.1630700466; nmstat=b566e578-9693-26ec-d4e0-67715efc3953; scarab.visitor=%2234D44C436B0EA2B6%22; __pr.1214=fuSYf87ma0; dtSa=true%7CKU%7C-1%7CPage%3A%20MCRZRV1-35865.html%7C-%7C1630700468081%7C500465985_432%7Chttps%3A%2F%2Fwww.newbalance.com%2Fpd%2Ffresh-foam-cruzv1-reissue%2FMCRZRV1-35865.html%7CFresh%20Foam%20Cruzv1%20Reissue%20-%20New%20Balance%7C1630700466839%7C%7C; rxVisitor=16307004659907GAV627UN0SU7LG34O5QKUDIBVB77QVT; __cq_seg=0~-0.07!1~0.05!2~0.29!3~-0.78!4~-0.20!5~-0.34!6~-0.06!7~0.21!8~0.28!9~0.15; _uetsid=7b544f600cf411ec837e1b9e7f1769f7; _uetvid=7b5467e00cf411ec922a6b6fa3f45bb6; rmStore=amid:39756; scarab.profile=%22MCRZRV1%252D35865%7C1630700472%22; cto_bundle=_HO4319vSmh2NlMwWjhYc1NHZ2dFWUxBZkdiMm9PUUd0T2xGYVRBc3JWcU15ZE9qWXRuNzl3VVVrQjVaQnNyZ0V2Vm1xV1pHUmxWbDYzbk4yb0UzTTZsZk5aYWdYMyUyQmtrV0RQUGs0WmtlMG1yJTJGaCUyRjg0NXF0WmtSY2poRHZYMmhLaEtkbA; _fbp=fb.1.1630700472684.1539229442; OptanonConsent=isIABGlobal=false&datestamp=Fri+Sep+03+2021+16%3A21%3A12+GMT-0400+(Eastern+Daylight+Time)&version=6.19.0&hosts=&landingPath=https%3A%2F%2Fwww.newbalance.com%2Fpd%2Ffresh-foam-cruzv1-reissue%2FMCRZRV1-35865.html&groups=C0001%3A1%2CC0003%3A1%2CBG8%3A1%2CC0002%3A1%2CC0004%3A1; stc112222=tsa:1630700472805.1582937491.8139043.5932129214534736.:20210903205112|env:1%7C20211004202112%7C20210903205112%7C1%7C1022688:20220903202112|uid:1630700472804.1140644789.4314075.112222.1549793015:20220903202112|srchist:1022688%3A1%3A20211004202112:20220903202112; hideEmail=hide; utag_vs=3; utag_dslv=1630700482186; RT="z=1&dm=newbalance.com&si=53d7ef34-4132-405e-812b-bf8e0d095df9&ss=kt4sybok&sl=2&tt=2k4&bcn=%2F%2F173c5b05.akstat.io%2F&r=35df020293b122bf6aac271e3c596a40&obo=1&hd=dhv"; _cs_id=04ff01a6-af87-a24d-9860-261b64bcab6a.1630700465.2.1630702292.1630702292.1.1664864465322; _cs_s=1.5.0.1630704092292; RT="z=1&dm=www.newbalance.com&si=53d7ef34-4132-405e-812b-bf8e0d095df9&ss=kt4sybok&sl=2&tt=2k4&bcn=%2F%2F173c5b05.akstat.io%2F&obo=1"; utag_ppv=-,18,11,1306; _gat_tealium_0=1; utag_main=v_id:017bad526cb9001a4a34b3e43dd903073004a06b00fda$_sn:2$_se:2$_ss:0$_st:1630705678495$_prevpage:Fresh%20Foam%20Cruzv1%20Reissue%3Bexp-1630707478498$dc_visit:2$vapi_domain:newbalance.com$ses_id:1630703877899%3Bexp-session$_pn:1%3Bexp-session$dc_event:2%3Bexp-session$dc_region:us-east-1%3Bexp-session; dtLatC=9; dtCookie=v_4_srv_2_sn_KF7U154126CN4US83L8ITJS758JP1KSF_app-3A394e822a8075a68e_0_ol_0_perc_100000_mul_1; s_sq=nbusalive%3D%2526c.%2526a.%2526activitymap.%2526page%253DFresh%252520Foam%252520Cruzv1%252520Reissue%2526link%253DAdd%252520to%252520cart%2526region%253DproductAtributes%2526pageIDType%253D1%2526.activitymap%2526.a%2526.c%2526pid%253DFresh%252520Foam%252520Cruzv1%252520Reissue%2526pidt%253D1%2526oid%253DAdd%252520to%252520cart%2526oidt%253D3%2526ot%253DSUBMIT; _gali=productAtributes; rxvt=1630705698787|1630702273996; dtPC=2$500472369_531h28vPUNGQCPHUHIIUBHMKQIJELBUOKKUAIFM-0e2`},
		http.PseudoHeaderOrderKey: {http.PseudoMethod, http.PseudoAuthority, http.PseudoScheme, http.PseudoPath},
		http.HeaderOrderKey:
			{
			"Host",
			"Connection",
			"Content-Length",
			"Pragma",
			"Cache-Control",
			"sec-ch-ua",
			"PRMeACzwub-f",
			"PRMeACzwub-a",
			"sec-ch-ua-mobile",
			"User-Agent",
			"PRMeACzwub-c",
			"PRMeACzwub-z",
			"Accept",
			"PRMeACzwub-b",
			"x-dtpc",
			"X-Requested-With",
			"Content-Type",
			"PRMeACzwub-d",
			"Origin",
			"Sec-Fetch-Site",
			"Sec-Fetch-Mode",
			"Sec-Fetch-Dest",
			"Referer",
			"Accept-Encoding",
			"Accept-Language",
			"Cookie",
			},
	}


	tk.SetStatus(module.STATUS_ADDING_TO_CART, "adding to cart")
	tk.Do(req)
	tk.SetStatus(module.STATUS_WAITING_FOR_CHECKOUT, "moving to checkout")
}
