using System;
using System.Collections.Generic;
using System.Linq;
using System.Reflection;
using System.Runtime.InteropServices;

using ILRuntime.CLR.TypeSystem;
using ILRuntime.CLR.Method;
using ILRuntime.Runtime.Enviorment;
using ILRuntime.Runtime.Intepreter;
using ILRuntime.Runtime.Stack;
using ILRuntime.Reflection;
using ILRuntime.CLR.Utils;
#if DEBUG && !DISABLE_ILRUNTIME_DEBUG
using AutoList = System.Collections.Generic.List<object>;
#else
using AutoList = ILRuntime.Other.UncheckedList<object>;
#endif
namespace ILRuntime.Runtime.Generated
{
    unsafe class JEngine_Core_Updater_Binding
    {
        public static void Register(ILRuntime.Runtime.Enviorment.AppDomain app)
        {
            BindingFlags flag = BindingFlags.Public | BindingFlags.Instance | BindingFlags.Static | BindingFlags.DeclaredOnly;
            MethodBase method;
            Type[] args;
            Type type = typeof(JEngine.Core.Updater);
            args = new Type[]{};
            method = type.GetMethod("Quit", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, Quit_0);
            args = new Type[]{};
            method = type.GetMethod("get_Mode", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, get_Mode_1);
            args = new Type[]{};
            method = type.GetMethod("get_WebVersion", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, get_WebVersion_2);
            args = new Type[]{};
            method = type.GetMethod("get_ResourceUrl", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, get_ResourceUrl_3);
            args = new Type[]{};
            method = type.GetMethod("get_MainPackageName", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, get_MainPackageName_4);


        }


        static StackObject* Quit_0(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* __ret = ILIntepreter.Minus(__esp, 0);


            JEngine.Core.Updater.Quit();

            return __ret;
        }

        static StackObject* get_Mode_1(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* __ret = ILIntepreter.Minus(__esp, 0);


            var result_of_this_method = JEngine.Core.Updater.Mode;

            return ILIntepreter.PushObject(__ret, __mStack, result_of_this_method);
        }

        static StackObject* get_WebVersion_2(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* __ret = ILIntepreter.Minus(__esp, 0);


            var result_of_this_method = JEngine.Core.Updater.WebVersion;

            return ILIntepreter.PushObject(__ret, __mStack, result_of_this_method);
        }

        static StackObject* get_ResourceUrl_3(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* __ret = ILIntepreter.Minus(__esp, 0);


            var result_of_this_method = JEngine.Core.Updater.ResourceUrl;

            return ILIntepreter.PushObject(__ret, __mStack, result_of_this_method);
        }

        static StackObject* get_MainPackageName_4(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* __ret = ILIntepreter.Minus(__esp, 0);


            var result_of_this_method = JEngine.Core.Updater.MainPackageName;

            return ILIntepreter.PushObject(__ret, __mStack, result_of_this_method);
        }



    }
}
