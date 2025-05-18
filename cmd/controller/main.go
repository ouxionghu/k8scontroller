package main

// 导入所需的标准库和第三方库
import (
	"flag" // 用于解析命令行参数
	"os"   // 用于操作系统交互，如退出程序

	"k8s.io/apimachinery/pkg/runtime"                   // Kubernetes API 对象的运行时工具包
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"  // 运行时工具包，包含错误处理等
	clientgoscheme "k8s.io/client-go/kubernetes/scheme" // 用于注册 Kubernetes 内置资源类型到 scheme
	ctrl "sigs.k8s.io/controller-runtime"               // controller-runtime 的主入口，提供 manager、controller 等
	"sigs.k8s.io/controller-runtime/pkg/healthz"        // 健康检查相关工具
	"sigs.k8s.io/controller-runtime/pkg/log/zap"        // zap 日志库的适配器

	"github.com/ouxionghu/k8scontroller/pkg/controller" // 引入自定义 controller 逻辑
)

// 定义全局变量
var (
	scheme   = runtime.NewScheme()        // scheme 用于注册所有资源类型（内置和自定义）
	setupLog = ctrl.Log.WithName("setup") // setupLog 用于记录启动和初始化相关日志
)

// init 函数在 main 之前自动执行，用于初始化 scheme
func init() {
	// 注册 Kubernetes 内置资源类型到 scheme，确保 manager 能识别这些类型
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	// 如果有自定义资源类型，也需要在这里注册
	// 例如：utilruntime.Must(examplev1.AddToScheme(scheme))
}

// main 函数是程序入口
func main() {
	var metricsAddr string        // 用于存储 metrics 监听地址
	var enableLeaderElection bool // 是否启用 leader 选举
	var probeAddr string          // 用于存储健康检查监听地址

	// 解析命令行参数，允许用户自定义 metrics 和 probe 的监听地址，以及是否启用 leader 选举
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")

	// 配置 zap 日志库的开发模式选项
	opts := zap.Options{
		Development: true, // 开发模式下日志更友好
	}
	opts.BindFlags(flag.CommandLine) // 允许通过命令行参数配置日志
	flag.Parse()                     // 解析所有命令行参数

	// 设置全局 logger，后续 controller-runtime 都会用这个 logger
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	// 创建 manager，负责管理 controller、webhook、metrics、健康检查等
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,                                  // 注册的资源类型集合
		HealthProbeBindAddress: probeAddr,                               // 健康检查监听地址
		LeaderElection:         enableLeaderElection,                    // 是否启用 leader 选举
		LeaderElectionID:       "k8scontroller.yourusername.github.com", // leader 选举的唯一标识
	})
	// 如果 manager 创建失败，记录错误并退出程序
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// 注册自定义 controller 到 manager
	if err = (&controller.Controller{
		Client: mgr.GetClient(), // 注入 client 用于与 Kubernetes API 交互
		Scheme: mgr.GetScheme(), // 注入 scheme 用于资源类型识别
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller") // 注册失败则记录错误并退出
		os.Exit(1)
	}

	// 添加健康检查端点，Kubernetes 可通过 /healthz 检查 controller 健康
	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	// 添加就绪检查端点，Kubernetes 可通过 /readyz 检查 controller 是否就绪
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	// 启动 manager，开始监听和处理事件（阻塞主线程，直到收到终止信号）
	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
