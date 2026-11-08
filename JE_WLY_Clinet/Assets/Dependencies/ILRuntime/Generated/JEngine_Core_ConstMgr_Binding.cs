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
    unsafe class JEngine_Core_ConstMgr_Binding
    {
        public static void Register(ILRuntime.Runtime.Enviorment.AppDomain app)
        {
            BindingFlags flag = BindingFlags.Public | BindingFlags.Instance | BindingFlags.Static | BindingFlags.DeclaredOnly;
            FieldInfo field;
            Type[] args;
            Type type = typeof(JEngine.Core.ConstMgr);

            field = type.GetField("WaitFor1Sec", flag);
            app.RegisterCLRFieldGetter(field, get_WaitFor1Sec_0);
            app.RegisterCLRFieldSetter(field, set_WaitFor1Sec_0);
            app.RegisterCLRFieldBinding(field, CopyToStack_WaitFor1Sec_0, AssignFromStack_WaitFor1Sec_0);


        }



        static object get_WaitFor1Sec_0(ref object o)
        {
            return JEngine.Core.ConstMgr.WaitFor1Sec;
        }

        static StackObject* CopyToStack_WaitFor1Sec_0(ref object o, ILIntepreter __intp, StackObject* __ret, AutoList __mStack)
        {
            var result_of_this_method = JEngine.Core.ConstMgr.WaitFor1Sec;
            object obj_result_of_this_method = result_of_this_method;
            if(obj_result_of_this_method is CrossBindingAdaptorType)
            {    
                return ILIntepreter.PushObject(__ret, __mStack, ((CrossBindingAdaptorType)obj_result_of_this_method).ILInstance);
            }
            return ILIntepreter.PushObject(__ret, __mStack, result_of_this_method);
        }

        static void set_WaitFor1Sec_0(ref object o, object v)
        {
            JEngine.Core.ConstMgr.WaitFor1Sec = (UnityEngine.WaitForSecondsRealtime)v;
        }

        static StackObject* AssignFromStack_WaitFor1Sec_0(ref object o, ILIntepreter __intp, StackObject* ptr_of_this_method, AutoList __mStack)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            UnityEngine.WaitForSecondsRealtime @WaitFor1Sec = (UnityEngine.WaitForSecondsRealtime)typeof(UnityEngine.WaitForSecondsRealtime).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            JEngine.Core.ConstMgr.WaitFor1Sec = @WaitFor1Sec;
            return ptr_of_this_method;
        }



    }
}
