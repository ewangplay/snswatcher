/* 
 * thrift interface for snswatcher
 */

namespace cpp jzlservice.snswatcher
namespace go jzlservice.snswatcher
namespace py jzlservice.snswatcher
namespace php jzlservice.snswatcher
namespace perl jzlservice.snswatcher
namespace java jzlservice.snswatcher

/**
* snswatcher service
*/
service SNSWatcher {
	string ping(),		            //服务的连通性测试接口
                                    //返回: pong
}

