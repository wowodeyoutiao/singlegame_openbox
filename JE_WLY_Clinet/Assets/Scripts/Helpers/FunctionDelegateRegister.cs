using JEngine.Core;
using AppDomain = ILRuntime.Runtime.Enviorment.AppDomain;

namespace JEngine.Helper
{
    public class FunctionDelegateRegister : IRegisterHelper
    {
        public void Register(AppDomain appdomain)
        {
            /*  这里注册你所需要的FunctionDelegate，参考报错的内容黏贴就好，举个例子：
             *
                appdomain.DelegateManager.RegisterFunctionDelegate<System.String, System.Boolean>();
             *
             */
            appdomain.DelegateManager.RegisterFunctionDelegate<System.Int32, System.Int32>();
            appdomain.DelegateManager.RegisterFunctionDelegate<System.Int32, System.Boolean>();
            appdomain.DelegateManager.RegisterFunctionDelegate<System.Int64, System.Boolean>();
            appdomain.DelegateManager.RegisterFunctionDelegate<System.Single, System.Boolean>();
            appdomain.DelegateManager.RegisterFunctionDelegate<System.Double, System.Boolean>();
            appdomain.DelegateManager.RegisterFunctionDelegate<System.Boolean, System.Boolean>();
            appdomain.DelegateManager.RegisterFunctionDelegate<System.Object, System.Boolean>();

            appdomain.DelegateManager.RegisterFunctionDelegate<System.Collections.Generic.Dictionary<System.Object, System.Collections.Generic.List<System.ValueTuple<System.Int32, System.Action<System.Object, System.Object>>>>, System.Boolean>();


            appdomain.DelegateManager.RegisterFunctionDelegate<Cysharp.Threading.Tasks.UniTaskVoid>();
            appdomain.DelegateManager.RegisterFunctionDelegate<ILRuntime.Runtime.Intepreter.ILTypeInstance, System.Int32>();
            appdomain.DelegateManager.RegisterFunctionDelegate<ILRuntime.Runtime.Intepreter.ILTypeInstance, System.Int64>();
            appdomain.DelegateManager.RegisterFunctionDelegate<ILRuntime.Runtime.Intepreter.ILTypeInstance, System.Single>();
            appdomain.DelegateManager.RegisterFunctionDelegate<ILRuntime.Runtime.Intepreter.ILTypeInstance, System.Double>();
            appdomain.DelegateManager.RegisterFunctionDelegate<ILRuntime.Runtime.Intepreter.ILTypeInstance, ILRuntime.Runtime.Intepreter.ILTypeInstance>();
            appdomain.DelegateManager.RegisterFunctionDelegate<ILRuntime.Runtime.Intepreter.ILTypeInstance, ILRuntime.Runtime.Intepreter.ILTypeInstance, System.Int32>();
            appdomain.DelegateManager.RegisterFunctionDelegate<System.Int32>();
            appdomain.DelegateManager.RegisterFunctionDelegate<System.Int32, System.Int32, System.Int32>();
            appdomain.DelegateManager.RegisterFunctionDelegate<SuperScrollView.LoopListView2, System.Int32, SuperScrollView.LoopListViewItem2>();
            appdomain.DelegateManager.RegisterFunctionDelegate<SuperScrollView.LoopGridView, System.Int32, System.Int32, System.Int32, SuperScrollView.LoopGridViewItem>();
            appdomain.DelegateManager.RegisterMethodDelegate<System.Int64, System.String>();

        }
    }
}