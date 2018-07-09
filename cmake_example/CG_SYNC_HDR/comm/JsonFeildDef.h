#ifndef __JSON_FEILD_DEF_H_
#define __JSON_FEILD_DEF_H_

const char* const Code                 = "A30"; //��Ʊ����     
const char* const UpdateDate           = "A99"; //���ڣ���ʽ��Unix��ʽʱ����ֵ��
const char* const Num                  = "A98"; //��������
const char* const Data                 = "A97"; //���������ֶ�
/* �����ʽ�����*/
const char* const AmountInRate         = "A01"; //������ռ�ɽ����                                                             
const char* const AreaAmountRate       = "A02"; //ռ�����ܳɽ�����                                                    
const char* const AreaChangeRate       = "A03"; //���任����                                                          
const char* const AreaClose            = "A04"; //�������̼�                                                          
const char* const AreaDate             = "A05"; //                                                                    
const char* const AreaDeltaPercent     = "A06"; //�����ǵ���                                                          
const char* const ChangeRate           = "A07"; //������                                                              
const char* const DayCount             = "A08"; //��������(0=���ա�1=3�ա�2=5�ա�3=10�ա�4=20��)                      
const char* const DeltaPercent         = "A09"; //�ǵ���                                                              
const char* const ImpetusBill          = "A10"; //�󵥶���                                                            
const char* const IndustryCode         = "A11"; //��ҵ����                                                            
const char* const MainNetFlowIn        = "A12"; //����������(Ԫ)                                                      
const char* const MainNetFlowOut       = "A13"; //����������(Ԫ)                                                      
const char* const MainNetInRate        = "A14"; //��������������                                                      
const char* const MainNetOutRate       = "A15"; //��������������                                                      
const char* const NetBigBill           = "A16"; //�󵥾���(Ԫ)                                                        
const char* const NetFlowIn            = "A17"; //������(Ԫ)                                                          
const char* const NetFlowInPower       = "A18"; //����������                                                          
const char* const NetFlowOut           = "A19"; //������(Ԫ)                                                          
const char* const NetFlowOutPower      = "A20"; //����������                                                          
const char* const NowPrice             = "A21"; //�ּ�                                                                
const char* const SeriesAddDays        = "A22"; //������������(SeriesNetIn>0��ʾ�������֣� SeriesNetOut>0��ʾ��������)
const char* const SeriesNetIn          = "A23"; //�������ֻ���������                                                  
const char* const SeriesNetInPower     = "A24"; //�������ֻ�������������                                              
const char* const SeriesNetOut         = "A25"; //�������ֻ���������                                                  
const char* const SeriesNetOutPower    = "A26"; //�������ֻ�������������                                              
const char* const TimeStamp            = "A27"; //ʱ���                                                              
const char* const TotalFlowIn          = "A28"; //������(Ԫ)                                                          
const char* const TotalFlowOut         = "A29"; //������(Ԫ)
const char* const nDayClose            = "A31"; //3�ջ���5�ջ����������ڵ����̼�                                      
const char* const nDayDate             = "A32"; //                                                                     

/* �����ʽ�������ϸ */
const char* const Date					= "B01"; //	���ڣ���ʽ��Unix��ʽʱ����ֵ��
const char* const BuyBigOrderVolume		= "B02"; //	�������
const char* const BuyBigOrderAmount		= "B03"; //	����󵥶�
const char* const SelBigOrderVolume		= "B04"; //	��������
const char* const SelBigOrderAmount		= "B05"; //	�����󵥶�
const char* const BuyMinOrderVolume		= "B06"; //	����С����
const char* const BuyMinOrderAmount		= "B07"; //	����С����
const char* const SelMinOrderVolume		= "B08"; //	����С����
const char* const SelMinOrderAmount		= "B09"; //	����С����
const char* const BuyMidOrderVolume		= "B10"; //	�����е���
const char* const BuyMidOrderAmount		= "B11"; //	�����е���
const char* const SelMidOrderVolume		= "B12"; //	�����е���
const char* const SelMidOrderAmount		= "B13"; //	�����е���
const char* const BuyLargeOrderVolume	= "B14"; //	�����ش���
const char* const BuyLargeOrderAmount	= "B15"; //	�����ش󵥶�
const char* const SelLargeOrderVolume	= "B16"; //	�����ش���
const char* const SelLargeOrderAmount	= "B17"; //	�����ش󵥶�
const char* const DealCountVolume		= "B18"; //	�ɽ�������
const char* const DealCountAmount		= "B19"; //	�ɽ���

/* ������-���Ʒ�������
 * 0x0A69
 * */
const char* const ShortSupportPrice			= "D2"; //����֧��λ                         
const char* const ShortResistancePrice      = "D3"; //��������λ                         
const char* const MidSupportPrice           = "D4"; //����֧��λ                         
const char* const MidResistancePrice        = "D5"; //��������λ                         
const char* const Glossary                  = "D6"; //��������(MACD����ָ��)��utf-8���룩
 

/*�ʽ���-��������-�����ɱ����� 
 * 0x0A70*/
const char* const ClosePrice	= "E3"; //���̼�                        
const char* const MarketPrice	= "E4"; //10���г�ƽ�����׳ɱ�          
const char* const MainPrice		= "E5"; //10������ƽ�����׳ɱ�    

/*������-�г��������� 
 * 0x0A71*/
const char* const AccumulativeMarkup		= "F2"; //20���ۼ��Ƿ�                  
const char* const AccumulativeMarkupHY		= "F3"; //20����ҵ�ۼ��Ƿ�              
const char* const AccumulativeMarkupSHZS	= "F4"; //20����ָ֤���ۼ��Ƿ�          
const char* const TurnOverRate1				= "F5"; //1�ջ�����                     
const char* const TurnOverRate3				= "F6"; //3�ջ�����                     
const char* const TurnOverRate5				= "F7"; //5�ջ�����                     

/*������-����������� 
 * 0x0A72*/
const char* const YYSR			= "H2"; //Ӫҵ����        
const char* const YYSR_TB_ZZ	= "H3";	//Ӫҵ����ͬ������
const char* const JLR_PARENT	= "H4"; //������          
const char* const JLR_TB_ZZ		= "H5"; //������ͬ������  
const char* const MGSY			= "H6"; //ÿ������        

/* ������-�������-����ָ������
 * 0x0A73*/
const char* const YYSR_TB_ZZL	= "I3"; //Ӫҵ����ͬ��������            
const char* const JLR_TB_ZZL	= "I4"; //������ͬ��������              

/* ������-�������-��ҵ�������� 
 * 0x0A74*/
const char* const IndustryName			= "K02"; //������ҵ����        
const char* const ZZC					= "K03"; //�������ʲ�          
const char* const ZZC_HY                = "K04"; //��ҵ���ʲ���ֵ      
const char* const ZZC_PM                = "K05"; //�������ʲ�����      
const char* const JLRL                  = "K06"; //���ɾ�������        
const char* const JLRL_HY               = "K07"; //��ҵ��������        
const char* const JLRL_PM               = "K08"; //���ɾ�����������    
const char* const JZCSYL                = "K09"; //���ɾ��ʲ�������    
const char* const JZCSYL_HY             = "K10"; //��ҵ���ʲ�������    
const char* const JZCSYL_PM             = "K11"; //���ɾ��ʲ�����������
const char* const XSMLL                 = "K12"; //��������ë����      
const char* const XSMLL_HY             = "K13"; //��ҵ����ë����      
const char* const XSMLL_PM              = "K14"; //��������ë��������  
const char* const InduElementNum        = "K15"; //������ҵ�ɷݹ�����  

/*������-��ֵ�о����� 
 * 0x0A75*/
const char* const PE_TTM			= "L02";	//��ӯ��    
const char* const PE_TTM_HY         = "L03";	//��ҵ��ӯ��
const char* const PE_TTM_MARKET     = "L04";	//�г���ӯ��
const char* const PE_TTM_RANKING    = "L05";	//��ӯ������
const char* const PB                = "L06";	//�о���    
const char* const PB_HY             = "L07";	//��ҵ�о���
const char* const PB_MARKET         = "L08";	//�г��о���
const char* const PB_RANKING        = "L09";	//�о�������
const char* const PEG               = "L10";	//PEG       
const char* const PEG_HY            = "L11";	//��ҵPEG   
const char* const PEG_MARKET        = "L12";	//�г�PEG   
const char* const PEG_RANKING       = "L13";	//PEG����   

/* ������-�����۵�����
 * 0x0A76*/
const char* const OrganNum      	= "M02"; //��������  
const char* const StatDays          = "M03"; //ͳ������  
const char* const Rating            = "M04"; //��������  
const char* const Attention         = "M05"; //������ע��
const char* const BuyNum            = "M06"; //��������  
const char* const HoldingsNum       = "M07"; //��������  
const char* const NeutralNum        = "M08"; //��������  
const char* const ReductionNum      = "M09"; //��������  
const char* const SellNum           = "M10"; //��������  


#endif

