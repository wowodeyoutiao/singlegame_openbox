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
    unsafe class Spine_Unity_SkeletonRenderer_Binding
    {
        public static void Register(ILRuntime.Runtime.Enviorment.AppDomain app)
        {
            BindingFlags flag = BindingFlags.Public | BindingFlags.Instance | BindingFlags.Static | BindingFlags.DeclaredOnly;
            FieldInfo field;
            Type[] args;
            Type type = typeof(Spine.Unity.SkeletonRenderer);

            field = type.GetField("skeleton", flag);
            app.RegisterCLRFieldGetter(field, get_skeleton_0);
            app.RegisterCLRFieldSetter(field, set_skeleton_0);
            app.RegisterCLRFieldBinding(field, CopyToStack_skeleton_0, AssignFromStack_skeleton_0);
            field = type.GetField("skeletonDataAsset", flag);
            app.RegisterCLRFieldGetter(field, get_skeletonDataAsset_1);
            app.RegisterCLRFieldSetter(field, set_skeletonDataAsset_1);
            app.RegisterCLRFieldBinding(field, CopyToStack_skeletonDataAsset_1, AssignFromStack_skeletonDataAsset_1);


        }



        static object get_skeleton_0(ref object o)
        {
            return ((Spine.Unity.SkeletonRenderer)o).skeleton;
        }

        static StackObject* CopyToStack_skeleton_0(ref object o, ILIntepreter __intp, StackObject* __ret, AutoList __mStack)
        {
            var result_of_this_method = ((Spine.Unity.SkeletonRenderer)o).skeleton;
            return ILIntepreter.PushObject(__ret, __mStack, result_of_this_method);
        }

        static void set_skeleton_0(ref object o, object v)
        {
            ((Spine.Unity.SkeletonRenderer)o).skeleton = (Spine.Skeleton)v;
        }

        static StackObject* AssignFromStack_skeleton_0(ref object o, ILIntepreter __intp, StackObject* ptr_of_this_method, AutoList __mStack)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            Spine.Skeleton @skeleton = (Spine.Skeleton)typeof(Spine.Skeleton).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            ((Spine.Unity.SkeletonRenderer)o).skeleton = @skeleton;
            return ptr_of_this_method;
        }

        static object get_skeletonDataAsset_1(ref object o)
        {
            return ((Spine.Unity.SkeletonRenderer)o).skeletonDataAsset;
        }

        static StackObject* CopyToStack_skeletonDataAsset_1(ref object o, ILIntepreter __intp, StackObject* __ret, AutoList __mStack)
        {
            var result_of_this_method = ((Spine.Unity.SkeletonRenderer)o).skeletonDataAsset;
            return ILIntepreter.PushObject(__ret, __mStack, result_of_this_method);
        }

        static void set_skeletonDataAsset_1(ref object o, object v)
        {
            ((Spine.Unity.SkeletonRenderer)o).skeletonDataAsset = (Spine.Unity.SkeletonDataAsset)v;
        }

        static StackObject* AssignFromStack_skeletonDataAsset_1(ref object o, ILIntepreter __intp, StackObject* ptr_of_this_method, AutoList __mStack)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            Spine.Unity.SkeletonDataAsset @skeletonDataAsset = (Spine.Unity.SkeletonDataAsset)typeof(Spine.Unity.SkeletonDataAsset).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            ((Spine.Unity.SkeletonRenderer)o).skeletonDataAsset = @skeletonDataAsset;
            return ptr_of_this_method;
        }



    }
}
