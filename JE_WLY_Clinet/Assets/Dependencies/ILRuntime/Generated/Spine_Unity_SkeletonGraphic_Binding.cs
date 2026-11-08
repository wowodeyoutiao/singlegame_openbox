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
    unsafe class Spine_Unity_SkeletonGraphic_Binding
    {
        public static void Register(ILRuntime.Runtime.Enviorment.AppDomain app)
        {
            BindingFlags flag = BindingFlags.Public | BindingFlags.Instance | BindingFlags.Static | BindingFlags.DeclaredOnly;
            MethodBase method;
            FieldInfo field;
            Type[] args;
            Type type = typeof(Spine.Unity.SkeletonGraphic);
            args = new Type[]{};
            method = type.GetMethod("get_AnimationState", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, get_AnimationState_0);
            args = new Type[]{};
            method = type.GetMethod("get_Skeleton", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, get_Skeleton_1);
            args = new Type[]{typeof(System.Boolean)};
            method = type.GetMethod("Initialize", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, Initialize_2);

            field = type.GetField("skeletonDataAsset", flag);
            app.RegisterCLRFieldGetter(field, get_skeletonDataAsset_0);
            app.RegisterCLRFieldSetter(field, set_skeletonDataAsset_0);
            app.RegisterCLRFieldBinding(field, CopyToStack_skeletonDataAsset_0, AssignFromStack_skeletonDataAsset_0);


        }


        static StackObject* get_AnimationState_0(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 1);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            Spine.Unity.SkeletonGraphic instance_of_this_method = (Spine.Unity.SkeletonGraphic)typeof(Spine.Unity.SkeletonGraphic).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);

            var result_of_this_method = instance_of_this_method.AnimationState;

            return ILIntepreter.PushObject(__ret, __mStack, result_of_this_method);
        }

        static StackObject* get_Skeleton_1(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 1);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            Spine.Unity.SkeletonGraphic instance_of_this_method = (Spine.Unity.SkeletonGraphic)typeof(Spine.Unity.SkeletonGraphic).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);

            var result_of_this_method = instance_of_this_method.Skeleton;

            return ILIntepreter.PushObject(__ret, __mStack, result_of_this_method);
        }

        static StackObject* Initialize_2(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 2);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            System.Boolean @overwrite = ptr_of_this_method->Value == 1;

            ptr_of_this_method = ILIntepreter.Minus(__esp, 2);
            Spine.Unity.SkeletonGraphic instance_of_this_method = (Spine.Unity.SkeletonGraphic)typeof(Spine.Unity.SkeletonGraphic).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);

            instance_of_this_method.Initialize(@overwrite);

            return __ret;
        }


        static object get_skeletonDataAsset_0(ref object o)
        {
            return ((Spine.Unity.SkeletonGraphic)o).skeletonDataAsset;
        }

        static StackObject* CopyToStack_skeletonDataAsset_0(ref object o, ILIntepreter __intp, StackObject* __ret, AutoList __mStack)
        {
            var result_of_this_method = ((Spine.Unity.SkeletonGraphic)o).skeletonDataAsset;
            return ILIntepreter.PushObject(__ret, __mStack, result_of_this_method);
        }

        static void set_skeletonDataAsset_0(ref object o, object v)
        {
            ((Spine.Unity.SkeletonGraphic)o).skeletonDataAsset = (Spine.Unity.SkeletonDataAsset)v;
        }

        static StackObject* AssignFromStack_skeletonDataAsset_0(ref object o, ILIntepreter __intp, StackObject* ptr_of_this_method, AutoList __mStack)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            Spine.Unity.SkeletonDataAsset @skeletonDataAsset = (Spine.Unity.SkeletonDataAsset)typeof(Spine.Unity.SkeletonDataAsset).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            ((Spine.Unity.SkeletonGraphic)o).skeletonDataAsset = @skeletonDataAsset;
            return ptr_of_this_method;
        }



    }
}
