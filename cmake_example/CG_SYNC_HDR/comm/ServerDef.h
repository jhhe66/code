#ifndef __SERVER_DEF_H_
#define __SERVER_DEF_H_

const unsigned int SHM_KEY = 0x10000000;

#if 0

const char MEMC_HOST[] = {"172.17.14.215:12001"};


const char HOST[] = {"183.62.112.116"};
const unsigned short PORT = 5033;

#endif 

#if 1
const char MEMC_HOST[] = {"172.30.8.239:12001"};


const char HOST[] = {"183.62.104.80"};
const unsigned short PORT = 5030;

#endif 
//const char KEY_FUNDS_PATH[] = ".FUNDSPATH";
//const char KEY_FUNDS_DETAIL[] = ".FUNDSDTL";

const unsigned short PACK_HEADER_LEN = 29;

enum MC_INDEX_KEY 
{
    MC_KEY_FUNDS_PATH = 0,
    MC_KEY_FUNDS_DTL,
    MC_KEY_TREND_MARKET,
    MC_KEY_TREND_ANYS,
    MC_KEY_FUND_FINANCE,
    MC_KEY_FUND_INDU_RANK,
    MC_KEY_FUND_PROFIT,
    MC_KEY_FUND_ALUA,
    MC_KEY_FUND_ORG,
    MC_KEY_ZHHS,
    MC_KEY_FIN_MAIN_COST
};

extern const char*  MEMC_KEY_LIST[MC_KEY_FIN_MAIN_COST + 1]; 

#if 0
const char*  const MEMC_KEY_LIST[MC_KEY_FIN_MAIN_COST + 1] = 
{ 
    ".FUNDS_PATH", 
    ".FUNDS_DTL",
    ".TREND_MARKET",
    ".TREND_ANYS",
    ".FUND_FINANCE",
    ".FUND_INDU_RANK",
    ".FUND_PROFIT",
    ".FUND_ALUA",
    ".FUND_ORG",
    ".ZHHS",
    ".FIN_MAIN_COST"
};

#endif


#endif
