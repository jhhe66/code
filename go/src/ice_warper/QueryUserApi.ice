module QueryUserApi{
    exception RequestCanceledException
    {
    };
    sequence<byte> bytes;
    interface QueryUserI{
        ["amd"] idempotent void GetUsersByApp(bytes req, out bytes res) throws RequestCanceledException;
        ["amd"] idempotent void GetUsersByRegids(bytes req, out bytes res) throws RequestCanceledException;
        //add new interface 
        ["amd"] idempotent void ValidUsers(bytes req, out bytes res) throws RequestCanceledException;
        ["amd"] idempotent void ValidRegids(bytes req, out bytes res) throws RequestCanceledException;
        ["amd"] idempotent void ValidAppkeys(bytes req, out bytes res) throws RequestCanceledException;
    };
};

