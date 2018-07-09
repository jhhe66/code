#ifndef __JSON_FEILD_DEF_H_
#define __JSON_FEILD_DEF_H_

const char* const Code                 = "A30"; //股票代码     
const char* const UpdateDate           = "A99"; //日期（格式：Unix格式时间数值）
const char* const Num                  = "A98"; //内容数量
const char* const Data                 = "A97"; //内容数组字段
/* 个股资金流向*/
const char* const AmountInRate         = "A01"; //净流入占成交额比                                                             
const char* const AreaAmountRate       = "A02"; //占区间总成交比例                                                    
const char* const AreaChangeRate       = "A03"; //区间换手率                                                          
const char* const AreaClose            = "A04"; //区间收盘价                                                          
const char* const AreaDate             = "A05"; //                                                                    
const char* const AreaDeltaPercent     = "A06"; //区间涨跌幅                                                          
const char* const ChangeRate           = "A07"; //换手率                                                              
const char* const DayCount             = "A08"; //日期类型(0=今日、1=3日、2=5日、3=10日、4=20日)                      
const char* const DeltaPercent         = "A09"; //涨跌幅                                                              
const char* const ImpetusBill          = "A10"; //大单动能                                                            
const char* const IndustryCode         = "A11"; //行业代码                                                            
const char* const MainNetFlowIn        = "A12"; //主力净流入(元)                                                      
const char* const MainNetFlowOut       = "A13"; //主力净流出(元)                                                      
const char* const MainNetInRate        = "A14"; //主力净流入力度                                                      
const char* const MainNetOutRate       = "A15"; //主力净流出力度                                                      
const char* const NetBigBill           = "A16"; //大单净量(元)                                                        
const char* const NetFlowIn            = "A17"; //净流入(元)                                                          
const char* const NetFlowInPower       = "A18"; //净流入力度                                                          
const char* const NetFlowOut           = "A19"; //总流出(元)                                                          
const char* const NetFlowOutPower      = "A20"; //净流入力度                                                          
const char* const NowPrice             = "A21"; //现价                                                                
const char* const SeriesAddDays        = "A22"; //连续增仓天数(SeriesNetIn>0表示连续增仓， SeriesNetOut>0表示连续减仓)
const char* const SeriesNetIn          = "A23"; //连续增仓机构净流入                                                  
const char* const SeriesNetInPower     = "A24"; //连续减仓机构净流入力度                                              
const char* const SeriesNetOut         = "A25"; //连续增仓机构净流出                                                  
const char* const SeriesNetOutPower    = "A26"; //连续减仓机构净流出力度                                              
const char* const TimeStamp            = "A27"; //时间戳                                                              
const char* const TotalFlowIn          = "A28"; //总流入(元)                                                          
const char* const TotalFlowOut         = "A29"; //总流出(元)
const char* const nDayClose            = "A31"; //3日或者5日或者其它日期的收盘价                                      
const char* const nDayDate             = "A32"; //                                                                     

/* 个股资金流向明细 */
const char* const Date					= "B01"; //	日期（格式：Unix格式时间数值）
const char* const BuyBigOrderVolume		= "B02"; //	买入大单量
const char* const BuyBigOrderAmount		= "B03"; //	买入大单额
const char* const SelBigOrderVolume		= "B04"; //	卖出大单量
const char* const SelBigOrderAmount		= "B05"; //	卖出大单额
const char* const BuyMinOrderVolume		= "B06"; //	买入小单量
const char* const BuyMinOrderAmount		= "B07"; //	买入小单额
const char* const SelMinOrderVolume		= "B08"; //	卖出小单量
const char* const SelMinOrderAmount		= "B09"; //	卖出小单额
const char* const BuyMidOrderVolume		= "B10"; //	买入中单量
const char* const BuyMidOrderAmount		= "B11"; //	买入中单额
const char* const SelMidOrderVolume		= "B12"; //	卖出中单量
const char* const SelMidOrderAmount		= "B13"; //	卖出中单额
const char* const BuyLargeOrderVolume	= "B14"; //	买入特大单量
const char* const BuyLargeOrderAmount	= "B15"; //	买入特大单额
const char* const SelLargeOrderVolume	= "B16"; //	卖出特大单量
const char* const SelLargeOrderAmount	= "B17"; //	卖出特大单额
const char* const DealCountVolume		= "B18"; //	成交笔数量
const char* const DealCountAmount		= "B19"; //	成交额

/* 技术面-趋势分析数据
 * 0x0A69
 * */
const char* const ShortSupportPrice			= "D2"; //短线支撑位                         
const char* const ShortResistancePrice      = "D3"; //短线阻力位                         
const char* const MidSupportPrice           = "D4"; //中线支撑位                         
const char* const MidResistancePrice        = "D5"; //中线阻力位                         
const char* const Glossary                  = "D6"; //分析术语(MACD技术指标)（utf-8编码）
 

/*资金面-主力分析-主力成本分析 
 * 0x0A70*/
const char* const ClosePrice	= "E3"; //收盘价                        
const char* const MarketPrice	= "E4"; //10日市场平均交易成本          
const char* const MainPrice		= "E5"; //10日主力平均交易成本    

/*技术面-市场表现数据 
 * 0x0A71*/
const char* const AccumulativeMarkup		= "F2"; //20日累计涨幅                  
const char* const AccumulativeMarkupHY		= "F3"; //20日行业累计涨幅              
const char* const AccumulativeMarkupSHZS	= "F4"; //20日上证指数累计涨幅          
const char* const TurnOverRate1				= "F5"; //1日换手率                     
const char* const TurnOverRate3				= "F6"; //3日换手率                     
const char* const TurnOverRate5				= "F7"; //5日换手率                     

/*基本面-财务分析数据 
 * 0x0A72*/
const char* const YYSR			= "H2"; //营业收入        
const char* const YYSR_TB_ZZ	= "H3";	//营业收入同比增长
const char* const JLR_PARENT	= "H4"; //净利润          
const char* const JLR_TB_ZZ		= "H5"; //净利润同比增长  
const char* const MGSY			= "H6"; //每股收益        

/* 基本面-财务分析-利润指标数据
 * 0x0A73*/
const char* const YYSR_TB_ZZL	= "I3"; //营业收入同比增长率            
const char* const JLR_TB_ZZL	= "I4"; //净利润同比增长率              

/* 基本面-财务分析-行业排名数据 
 * 0x0A74*/
const char* const IndustryName			= "K02"; //所属行业名称        
const char* const ZZC					= "K03"; //个股总资产          
const char* const ZZC_HY                = "K04"; //行业总资产均值      
const char* const ZZC_PM                = "K05"; //个股总资产排名      
const char* const JLRL                  = "K06"; //个股净利润率        
const char* const JLRL_HY               = "K07"; //行业净利润率        
const char* const JLRL_PM               = "K08"; //个股净利润率排名    
const char* const JZCSYL                = "K09"; //个股净资产收益率    
const char* const JZCSYL_HY             = "K10"; //行业净资产收益率    
const char* const JZCSYL_PM             = "K11"; //个股净资产收益率排名
const char* const XSMLL                 = "K12"; //个股销售毛利率      
const char* const XSMLL_HY             = "K13"; //行业销售毛利率      
const char* const XSMLL_PM              = "K14"; //个股销售毛利率排名  
const char* const InduElementNum        = "K15"; //所属行业成份股总数  

/*基本面-估值研究数据 
 * 0x0A75*/
const char* const PE_TTM			= "L02";	//市盈率    
const char* const PE_TTM_HY         = "L03";	//行业市盈率
const char* const PE_TTM_MARKET     = "L04";	//市场市盈率
const char* const PE_TTM_RANKING    = "L05";	//市盈率排名
const char* const PB                = "L06";	//市净率    
const char* const PB_HY             = "L07";	//行业市净率
const char* const PB_MARKET         = "L08";	//市场市净率
const char* const PB_RANKING        = "L09";	//市净率排名
const char* const PEG               = "L10";	//PEG       
const char* const PEG_HY            = "L11";	//行业PEG   
const char* const PEG_MARKET        = "L12";	//市场PEG   
const char* const PEG_RANKING       = "L13";	//PEG排名   

/* 基本面-机构观点数据
 * 0x0A76*/
const char* const OrganNum      	= "M02"; //机构数量  
const char* const StatDays          = "M03"; //统计天数  
const char* const Rating            = "M04"; //机构评级  
const char* const Attention         = "M05"; //机构关注度
const char* const BuyNum            = "M06"; //买入数量  
const char* const HoldingsNum       = "M07"; //增持数量  
const char* const NeutralNum        = "M08"; //中性数量  
const char* const ReductionNum      = "M09"; //减持数量  
const char* const SellNum           = "M10"; //卖出数量  


#endif

