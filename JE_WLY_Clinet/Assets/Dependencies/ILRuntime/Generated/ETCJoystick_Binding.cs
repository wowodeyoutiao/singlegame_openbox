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
    unsafe class ETCJoystick_Binding
    {
        public static void Register(ILRuntime.Runtime.Enviorment.AppDomain app)
        {
            BindingFlags flag = BindingFlags.Public | BindingFlags.Instance | BindingFlags.Static | BindingFlags.DeclaredOnly;
            FieldInfo field;
            Type[] args;
            Type type = typeof(global::ETCJoystick);

            field = type.GetField("onMove", flag);
            app.RegisterCLRFieldGetter(field, get_onMove_0);
            app.RegisterCLRFieldSetter(field, set_onMove_0);
            app.RegisterCLRFieldBinding(field, CopyToStack_onMove_0, AssignFromStack_onMove_0);
            field = type.GetField("onMoveStart", flag);
            app.RegisterCLRFieldGetter(field, get_onMoveStart_1);
            app.RegisterCLRFieldSetter(field, set_onMoveStart_1);
            app.RegisterCLRFieldBinding(field, CopyToStack_onMoveStart_1, AssignFromStack_onMoveStart_1);
            field = type.GetField("onMoveEnd", flag);
            app.RegisterCLRFieldGetter(field, get_onMoveEnd_2);
            app.RegisterCLRFieldSetter(field, set_onMoveEnd_2);
            app.RegisterCLRFieldBinding(field, CopyToStack_onMoveEnd_2, AssignFromStack_onMoveEnd_2);


        }



        static object get_onMove_0(ref object o)
        {
            return ((global::ETCJoystick)o).onMove;
        }

        static StackObject* CopyToStack_onMove_0(ref object o, ILIntepreter __intp, StackObject* __ret, AutoList __mStack)
        {
            var result_of_this_method = ((global::ETCJoystick)o).onMove;
            return ILIntepreter.PushObject(__ret, __mStack, result_of_this_method);
        }

        static void set_onMove_0(ref object o, object v)
        {
            ((global::ETCJoystick)o).onMove = (global::ETCJoystick.OnMoveHandler)v;
        }

        static StackObject* AssignFromStack_onMove_0(ref object o, ILIntepreter __intp, StackObject* ptr_of_this_method, AutoList __mStack)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            global::ETCJoystick.OnMoveHandler @onMove = (global::ETCJoystick.OnMoveHandler)typeof(global::ETCJoystick.OnMoveHandler).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            ((global::ETCJoystick)o).onMove = @onMove;
            return ptr_of_this_method;
        }

        static object get_onMoveStart_1(ref object o)
        {
            return ((global::ETCJoystick)o).onMoveStart;
        }

        static StackObject* CopyToStack_onMoveStart_1(ref object o, ILIntepreter __intp, StackObject* __ret, AutoList __mStack)
        {
            var result_of_this_method = ((global::ETCJoystick)o).onMoveStart;
            return ILIntepreter.PushObject(__ret, __mStack, result_of_this_method);
        }

        static void set_onMoveStart_1(ref object o, object v)
        {
            ((global::ETCJoystick)o).onMoveStart = (global::ETCJoystick.OnMoveStartHandler)v;
        }

        static StackObject* AssignFromStack_onMoveStart_1(ref object o, ILIntepreter __intp, StackObject* ptr_of_this_method, AutoList __mStack)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            global::ETCJoystick.OnMoveStartHandler @onMoveStart = (global::ETCJoystick.OnMoveStartHandler)typeof(global::ETCJoystick.OnMoveStartHandler).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            ((global::ETCJoystick)o).onMoveStart = @onMoveStart;
            return ptr_of_this_method;
        }

        static object get_onMoveEnd_2(ref object o)
        {
            return ((global::ETCJoystick)o).onMoveEnd;
        }

        static StackObject* CopyToStack_onMoveEnd_2(ref object o, ILIntepreter __intp, StackObject* __ret, AutoList __mStack)
        {
            var result_of_this_method = ((global::ETCJoystick)o).onMoveEnd;
            return ILIntepreter.PushObject(__ret, __mStack, result_of_this_method);
        }

        static void set_onMoveEnd_2(ref object o, object v)
        {
            ((global::ETCJoystick)o).onMoveEnd = (global::ETCJoystick.OnMoveEndHandler)v;
        }

        static StackObject* AssignFromStack_onMoveEnd_2(ref object o, ILIntepreter __intp, StackObject* ptr_of_this_method, AutoList __mStack)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            global::ETCJoystick.OnMoveEndHandler @onMoveEnd = (global::ETCJoystick.OnMoveEndHandler)typeof(global::ETCJoystick.OnMoveEndHandler).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            ((global::ETCJoystick)o).onMoveEnd = @onMoveEnd;
            return ptr_of_this_method;
        }



    }
}
