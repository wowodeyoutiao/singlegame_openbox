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
    unsafe class Spine_Unity_SkeletonAnimation_Binding
    {
        public static void Register(ILRuntime.Runtime.Enviorment.AppDomain app)
        {
            BindingFlags flag = BindingFlags.Public | BindingFlags.Instance | BindingFlags.Static | BindingFlags.DeclaredOnly;
            MethodBase method;
            FieldInfo field;
            Type[] args;
            Type type = typeof(Spine.Unity.SkeletonAnimation);
            args = new Type[]{};
            method = type.GetMethod("get_AnimationState", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, get_AnimationState_0);

            field = type.GetField("state", flag);
            app.RegisterCLRFieldGetter(field, get_state_0);
            app.RegisterCLRFieldSetter(field, set_state_0);
            app.RegisterCLRFieldBinding(field, CopyToStack_state_0, AssignFromStack_state_0);
            field = type.GetField("loop", flag);
            app.RegisterCLRFieldGetter(field, get_loop_1);
            app.RegisterCLRFieldSetter(field, set_loop_1);
            app.RegisterCLRFieldBinding(field, CopyToStack_loop_1, AssignFromStack_loop_1);


        }


        static StackObject* get_AnimationState_0(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 1);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            Spine.Unity.SkeletonAnimation instance_of_this_method = (Spine.Unity.SkeletonAnimation)typeof(Spine.Unity.SkeletonAnimation).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);

            var result_of_this_method = instance_of_this_method.AnimationState;

            return ILIntepreter.PushObject(__ret, __mStack, result_of_this_method);
        }


        static object get_state_0(ref object o)
        {
            return ((Spine.Unity.SkeletonAnimation)o).state;
        }

        static StackObject* CopyToStack_state_0(ref object o, ILIntepreter __intp, StackObject* __ret, AutoList __mStack)
        {
            var result_of_this_method = ((Spine.Unity.SkeletonAnimation)o).state;
            return ILIntepreter.PushObject(__ret, __mStack, result_of_this_method);
        }

        static void set_state_0(ref object o, object v)
        {
            ((Spine.Unity.SkeletonAnimation)o).state = (Spine.AnimationState)v;
        }

        static StackObject* AssignFromStack_state_0(ref object o, ILIntepreter __intp, StackObject* ptr_of_this_method, AutoList __mStack)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            Spine.AnimationState @state = (Spine.AnimationState)typeof(Spine.AnimationState).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            ((Spine.Unity.SkeletonAnimation)o).state = @state;
            return ptr_of_this_method;
        }

        static object get_loop_1(ref object o)
        {
            return ((Spine.Unity.SkeletonAnimation)o).loop;
        }

        static StackObject* CopyToStack_loop_1(ref object o, ILIntepreter __intp, StackObject* __ret, AutoList __mStack)
        {
            var result_of_this_method = ((Spine.Unity.SkeletonAnimation)o).loop;
            __ret->ObjectType = ObjectTypes.Integer;
            __ret->Value = result_of_this_method ? 1 : 0;
            return __ret + 1;
        }

        static void set_loop_1(ref object o, object v)
        {
            ((Spine.Unity.SkeletonAnimation)o).loop = (System.Boolean)v;
        }

        static StackObject* AssignFromStack_loop_1(ref object o, ILIntepreter __intp, StackObject* ptr_of_this_method, AutoList __mStack)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            System.Boolean @loop = ptr_of_this_method->Value == 1;
            ((Spine.Unity.SkeletonAnimation)o).loop = @loop;
            return ptr_of_this_method;
        }



    }
}
