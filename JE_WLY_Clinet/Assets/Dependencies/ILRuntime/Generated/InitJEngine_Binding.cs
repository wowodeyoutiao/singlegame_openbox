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
    unsafe class InitJEngine_Binding
    {
        public static void Register(ILRuntime.Runtime.Enviorment.AppDomain app)
        {
            BindingFlags flag = BindingFlags.Public | BindingFlags.Instance | BindingFlags.Static | BindingFlags.DeclaredOnly;
            FieldInfo field;
            Type[] args;
            Type type = typeof(global::InitJEngine);

            field = type.GetField("Instance", flag);
            app.RegisterCLRFieldGetter(field, get_Instance_0);
            app.RegisterCLRFieldSetter(field, set_Instance_0);
            app.RegisterCLRFieldBinding(field, CopyToStack_Instance_0, AssignFromStack_Instance_0);
            field = type.GetField("IsFreeLogin", flag);
            app.RegisterCLRFieldGetter(field, get_IsFreeLogin_1);
            app.RegisterCLRFieldSetter(field, set_IsFreeLogin_1);
            app.RegisterCLRFieldBinding(field, CopyToStack_IsFreeLogin_1, AssignFromStack_IsFreeLogin_1);
            field = type.GetField("key", flag);
            app.RegisterCLRFieldGetter(field, get_key_2);
            app.RegisterCLRFieldSetter(field, set_key_2);
            app.RegisterCLRFieldBinding(field, CopyToStack_key_2, AssignFromStack_key_2);
            field = type.GetField("debug", flag);
            app.RegisterCLRFieldGetter(field, get_debug_3);
            app.RegisterCLRFieldSetter(field, set_debug_3);
            app.RegisterCLRFieldBinding(field, CopyToStack_debug_3, AssignFromStack_debug_3);


        }



        static object get_Instance_0(ref object o)
        {
            return global::InitJEngine.Instance;
        }

        static StackObject* CopyToStack_Instance_0(ref object o, ILIntepreter __intp, StackObject* __ret, AutoList __mStack)
        {
            var result_of_this_method = global::InitJEngine.Instance;
            object obj_result_of_this_method = result_of_this_method;
            if(obj_result_of_this_method is CrossBindingAdaptorType)
            {    
                return ILIntepreter.PushObject(__ret, __mStack, ((CrossBindingAdaptorType)obj_result_of_this_method).ILInstance);
            }
            return ILIntepreter.PushObject(__ret, __mStack, result_of_this_method);
        }

        static void set_Instance_0(ref object o, object v)
        {
            global::InitJEngine.Instance = (global::InitJEngine)v;
        }

        static StackObject* AssignFromStack_Instance_0(ref object o, ILIntepreter __intp, StackObject* ptr_of_this_method, AutoList __mStack)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            global::InitJEngine @Instance = (global::InitJEngine)typeof(global::InitJEngine).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            global::InitJEngine.Instance = @Instance;
            return ptr_of_this_method;
        }

        static object get_IsFreeLogin_1(ref object o)
        {
            return ((global::InitJEngine)o).IsFreeLogin;
        }

        static StackObject* CopyToStack_IsFreeLogin_1(ref object o, ILIntepreter __intp, StackObject* __ret, AutoList __mStack)
        {
            var result_of_this_method = ((global::InitJEngine)o).IsFreeLogin;
            __ret->ObjectType = ObjectTypes.Integer;
            __ret->Value = result_of_this_method ? 1 : 0;
            return __ret + 1;
        }

        static void set_IsFreeLogin_1(ref object o, object v)
        {
            ((global::InitJEngine)o).IsFreeLogin = (System.Boolean)v;
        }

        static StackObject* AssignFromStack_IsFreeLogin_1(ref object o, ILIntepreter __intp, StackObject* ptr_of_this_method, AutoList __mStack)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            System.Boolean @IsFreeLogin = ptr_of_this_method->Value == 1;
            ((global::InitJEngine)o).IsFreeLogin = @IsFreeLogin;
            return ptr_of_this_method;
        }

        static object get_key_2(ref object o)
        {
            return ((global::InitJEngine)o).key;
        }

        static StackObject* CopyToStack_key_2(ref object o, ILIntepreter __intp, StackObject* __ret, AutoList __mStack)
        {
            var result_of_this_method = ((global::InitJEngine)o).key;
            return ILIntepreter.PushObject(__ret, __mStack, result_of_this_method);
        }

        static void set_key_2(ref object o, object v)
        {
            ((global::InitJEngine)o).key = (System.String)v;
        }

        static StackObject* AssignFromStack_key_2(ref object o, ILIntepreter __intp, StackObject* ptr_of_this_method, AutoList __mStack)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            System.String @key = (System.String)typeof(System.String).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            ((global::InitJEngine)o).key = @key;
            return ptr_of_this_method;
        }

        static object get_debug_3(ref object o)
        {
            return ((global::InitJEngine)o).debug;
        }

        static StackObject* CopyToStack_debug_3(ref object o, ILIntepreter __intp, StackObject* __ret, AutoList __mStack)
        {
            var result_of_this_method = ((global::InitJEngine)o).debug;
            __ret->ObjectType = ObjectTypes.Integer;
            __ret->Value = result_of_this_method ? 1 : 0;
            return __ret + 1;
        }

        static void set_debug_3(ref object o, object v)
        {
            ((global::InitJEngine)o).debug = (System.Boolean)v;
        }

        static StackObject* AssignFromStack_debug_3(ref object o, ILIntepreter __intp, StackObject* ptr_of_this_method, AutoList __mStack)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            System.Boolean @debug = ptr_of_this_method->Value == 1;
            ((global::InitJEngine)o).debug = @debug;
            return ptr_of_this_method;
        }



    }
}
