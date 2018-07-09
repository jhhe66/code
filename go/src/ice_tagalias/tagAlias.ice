module TagAlias{
    exception RequestCanceledException
    {
    };
    sequence<byte> bytes;
    interface TagAliasOp{
        //["ami","amd"] idempotent void request(string req, out string res) throws RequestCanceledException;
        ["ami","amd"] void request(string req, out string res) throws RequestCanceledException;
        // req next page data 
        ["ami","amd"] void request_next(string req, out string res) throws RequestCanceledException;
        
        ["amd"] idempotent void getTaglist(bytes req, out bytes res) throws RequestCanceledException;
        ["amd"] idempotent void getUsersByTag(bytes req, out bytes res) throws RequestCanceledException;
        ["amd"] idempotent void getUsersByAlias(bytes req, out bytes res) throws RequestCanceledException;
        ["amd"] idempotent void getTagsByUser(bytes req, out bytes res) throws RequestCanceledException;
        ["amd"] idempotent void getAliasByUser(bytes req, out bytes res) throws RequestCanceledException;
        ["amd"] idempotent void checkUserBelongTag(bytes req, out bytes res) throws RequestCanceledException;
        ["amd"] idempotent void checkUserBelongAlias(bytes req, out bytes res) throws RequestCanceledException;
        ["amd"] idempotent void getUsersCountByTag(bytes req, out bytes res) throws RequestCanceledException;
        ["amd"] idempotent void getUsersCountByAlias(bytes req, out bytes res) throws RequestCanceledException;
        ["amd"] idempotent void validateTags(bytes req, out bytes res) throws RequestCanceledException;
        ["amd"] idempotent void validateAlias(bytes req, out bytes res) throws RequestCanceledException;
        void shutdown();
    };
    interface TagAliasOpAdd extends TagAliasOp{
        ["amd"] idempotent void getTagsCountByUser(bytes req, out bytes res) throws RequestCanceledException;
    };
};
