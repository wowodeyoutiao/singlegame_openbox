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
    unsafe class TDAnalyticsUtil_Binding
    {
        public static void Register(ILRuntime.Runtime.Enviorment.AppDomain app)
        {
            BindingFlags flag = BindingFlags.Public | BindingFlags.Instance | BindingFlags.Static | BindingFlags.DeclaredOnly;
            MethodBase method;
            Type[] args;
            Type type = typeof(global::TDAnalyticsUtil);
            args = new Type[]{typeof(System.String), typeof(System.Collections.Generic.Dictionary<System.String, System.Object>), typeof(System.Boolean)};
            method = type.GetMethod("Track", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, Track_0);
            args = new Type[]{typeof(System.String)};
            method = type.GetMethod("Login", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, Login_1);
            args = new Type[]{typeof(System.String)};
            method = type.GetMethod("SetDistinctId", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, SetDistinctId_2);
            args = new Type[]{typeof(System.Collections.Generic.Dictionary<System.String, System.Object>)};
            method = type.GetMethod("SetSuperProperties", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, SetSuperProperties_3);


        }


        static StackObject* Track_0(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 3);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            System.Boolean @needFlush = ptr_of_this_method->Value == 1;

            ptr_of_this_method = ILIntepreter.Minus(__esp, 2);
            System.Collections.Generic.Dictionary<System.String, System.Object> @properties = (System.Collections.Generic.Dictionary<System.String, System.Object>)typeof(System.Collections.Generic.Dictionary<System.String, System.Object>).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 3);
            System.String @eventName = (System.String)typeof(System.String).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);


            global::TDAnalyticsUtil.Track(@eventName, @properties, @needFlush);

            return __ret;
        }

        static StackObject* Login_1(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 1);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            System.String @account = (System.String)typeof(System.String).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);


            global::TDAnalyticsUtil.Login(@account);

            return __ret;
        }

        static StackObject* SetDistinctId_2(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 1);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            System.String @distinctId = (System.String)typeof(System.String).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);


            global::TDAnalyticsUtil.SetDistinctId(@distinctId);

            return __ret;
        }

        static StackObject* SetSuperProperties_3(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 1);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            System.Collections.Generic.Dictionary<System.String, System.Object> @properties = (System.Collections.Generic.Dictionary<System.String, System.Object>)typeof(System.Collections.Generic.Dictionary<System.String, System.Object>).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);


            global::TDAnalyticsUtil.SetSuperProperties(@properties);

            return __ret;
        }



    }
}
