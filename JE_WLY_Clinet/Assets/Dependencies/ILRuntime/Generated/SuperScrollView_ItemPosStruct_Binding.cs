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
    unsafe class SuperScrollView_ItemPosStruct_Binding
    {
        public static void Register(ILRuntime.Runtime.Enviorment.AppDomain app)
        {
            BindingFlags flag = BindingFlags.Public | BindingFlags.Instance | BindingFlags.Static | BindingFlags.DeclaredOnly;
            FieldInfo field;
            Type[] args;
            Type type = typeof(SuperScrollView.ItemPosStruct);

            field = type.GetField("mItemIndex", flag);
            app.RegisterCLRFieldGetter(field, get_mItemIndex_0);
            app.RegisterCLRFieldSetter(field, set_mItemIndex_0);
            app.RegisterCLRFieldBinding(field, CopyToStack_mItemIndex_0, AssignFromStack_mItemIndex_0);
            field = type.GetField("mItemOffset", flag);
            app.RegisterCLRFieldGetter(field, get_mItemOffset_1);
            app.RegisterCLRFieldSetter(field, set_mItemOffset_1);
            app.RegisterCLRFieldBinding(field, CopyToStack_mItemOffset_1, AssignFromStack_mItemOffset_1);

            app.RegisterCLRCreateDefaultInstance(type, () => new SuperScrollView.ItemPosStruct());


        }

        static void WriteBackInstance(ILRuntime.Runtime.Enviorment.AppDomain __domain, StackObject* ptr_of_this_method, AutoList __mStack, ref SuperScrollView.ItemPosStruct instance_of_this_method)
        {
            ptr_of_this_method = ILIntepreter.GetObjectAndResolveReference(ptr_of_this_method);
            switch(ptr_of_this_method->ObjectType)
            {
                case ObjectTypes.Object:
                    {
                        __mStack[ptr_of_this_method->Value] = instance_of_this_method;
                    }
                    break;
                case ObjectTypes.FieldReference:
                    {
                        var ___obj = __mStack[ptr_of_this_method->Value];
                        if(___obj is ILTypeInstance)
                        {
                            ((ILTypeInstance)___obj)[ptr_of_this_method->ValueLow] = instance_of_this_method;
                        }
                        else
                        {
                            var t = __domain.GetType(___obj.GetType()) as CLRType;
                            t.SetFieldValue(ptr_of_this_method->ValueLow, ref ___obj, instance_of_this_method);
                        }
                    }
                    break;
                case ObjectTypes.StaticFieldReference:
                    {
                        var t = __domain.GetType(ptr_of_this_method->Value);
                        if(t is ILType)
                        {
                            ((ILType)t).StaticInstance[ptr_of_this_method->ValueLow] = instance_of_this_method;
                        }
                        else
                        {
                            ((CLRType)t).SetStaticFieldValue(ptr_of_this_method->ValueLow, instance_of_this_method);
                        }
                    }
                    break;
                 case ObjectTypes.ArrayReference:
                    {
                        var instance_of_arrayReference = __mStack[ptr_of_this_method->Value] as SuperScrollView.ItemPosStruct[];
                        instance_of_arrayReference[ptr_of_this_method->ValueLow] = instance_of_this_method;
                    }
                    break;
            }
        }


        static object get_mItemIndex_0(ref object o)
        {
            return ((SuperScrollView.ItemPosStruct)o).mItemIndex;
        }

        static StackObject* CopyToStack_mItemIndex_0(ref object o, ILIntepreter __intp, StackObject* __ret, AutoList __mStack)
        {
            var result_of_this_method = ((SuperScrollView.ItemPosStruct)o).mItemIndex;
            __ret->ObjectType = ObjectTypes.Integer;
            __ret->Value = result_of_this_method;
            return __ret + 1;
        }

        static void set_mItemIndex_0(ref object o, object v)
        {
            SuperScrollView.ItemPosStruct ins =(SuperScrollView.ItemPosStruct)o;
            ins.mItemIndex = (System.Int32)v;
            o = ins;
        }

        static StackObject* AssignFromStack_mItemIndex_0(ref object o, ILIntepreter __intp, StackObject* ptr_of_this_method, AutoList __mStack)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            System.Int32 @mItemIndex = ptr_of_this_method->Value;
            SuperScrollView.ItemPosStruct ins =(SuperScrollView.ItemPosStruct)o;
            ins.mItemIndex = @mItemIndex;
            o = ins;
            return ptr_of_this_method;
        }

        static object get_mItemOffset_1(ref object o)
        {
            return ((SuperScrollView.ItemPosStruct)o).mItemOffset;
        }

        static StackObject* CopyToStack_mItemOffset_1(ref object o, ILIntepreter __intp, StackObject* __ret, AutoList __mStack)
        {
            var result_of_this_method = ((SuperScrollView.ItemPosStruct)o).mItemOffset;
            __ret->ObjectType = ObjectTypes.Float;
            *(float*)&__ret->Value = result_of_this_method;
            return __ret + 1;
        }

        static void set_mItemOffset_1(ref object o, object v)
        {
            SuperScrollView.ItemPosStruct ins =(SuperScrollView.ItemPosStruct)o;
            ins.mItemOffset = (System.Single)v;
            o = ins;
        }

        static StackObject* AssignFromStack_mItemOffset_1(ref object o, ILIntepreter __intp, StackObject* ptr_of_this_method, AutoList __mStack)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            System.Single @mItemOffset = *(float*)&ptr_of_this_method->Value;
            SuperScrollView.ItemPosStruct ins =(SuperScrollView.ItemPosStruct)o;
            ins.mItemOffset = @mItemOffset;
            o = ins;
            return ptr_of_this_method;
        }



    }
}
