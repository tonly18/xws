package ziface

type IRouter interface {
	PreHandle(request IRequest) error  //在处理conn业务之前的钩子方法
	Handle(request IRequest) error     //处理conn业务的方法
	PostHandle(request IRequest) error //处理conn业务之后的钩子方法
}
