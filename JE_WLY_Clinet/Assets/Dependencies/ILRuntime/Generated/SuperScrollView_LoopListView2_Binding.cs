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
    unsafe class SuperScrollView_LoopListView2_Binding
    {
        public static void Register(ILRuntime.Runtime.Enviorment.AppDomain app)
        {
            BindingFlags flag = BindingFlags.Public | BindingFlags.Instance | BindingFlags.Static | BindingFlags.DeclaredOnly;
            MethodBase method;
            Type[] args;
            Type type = typeof(SuperScrollView.LoopListView2);
            args = new Type[]{typeof(System.Int32), typeof(System.Func<SuperScrollView.LoopListView2, System.Int32, SuperScrollView.LoopListViewItem2>), typeof(SuperScrollView.LoopListViewInitParam)};
            method = type.GetMethod("InitListView", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, InitListView_0);
            args = new Type[]{typeof(System.String)};
            method = type.GetMethod("NewListViewItem", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, NewListViewItem_1);
            args = new Type[]{typeof(System.Int32), typeof(System.Boolean)};
            method = type.GetMethod("SetListItemCount", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, SetListItemCount_2);
            args = new Type[]{};
            method = type.GetMethod("RefreshAllShownItem", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, RefreshAllShownItem_3);
            args = new Type[]{};
            method = type.GetMethod("get_ItemTotalCount", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, get_ItemTotalCount_4);
            args = new Type[]{typeof(System.Int32), typeof(System.Single)};
            method = type.GetMethod("MovePanelToItemIndex", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, MovePanelToItemIndex_5);
            args = new Type[]{typeof(System.Int32)};
            method = type.GetMethod("GetShownItemByItemIndex", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, GetShownItemByItemIndex_6);
            args = new Type[]{};
            method = type.GetMethod("GetFirstShownItemIndexAndOffset", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, GetFirstShownItemIndexAndOffset_7);
            args = new Type[]{};
            method = type.GetMethod("get_ArrangeType", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, get_ArrangeType_8);


        }


        static StackObject* InitListView_0(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 4);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            SuperScrollView.LoopListViewInitParam @initParam = (SuperScrollView.LoopListViewInitParam)typeof(SuperScrollView.LoopListViewInitParam).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 2);
            System.Func<SuperScrollView.LoopListView2, System.Int32, SuperScrollView.LoopListViewItem2> @onGetItemByIndex = (System.Func<SuperScrollView.LoopListView2, System.Int32, SuperScrollView.LoopListViewItem2>)typeof(System.Func<SuperScrollView.LoopListView2, System.Int32, SuperScrollView.LoopListViewItem2>).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)8);
            __intp.Free(ptr_of_this_method);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 3);
            System.Int32 @itemTotalCount = ptr_of_this_method->Value;

            ptr_of_this_method = ILIntepreter.Minus(__esp, 4);
            SuperScrollView.LoopListView2 instance_of_this_method = (SuperScrollView.LoopListView2)typeof(SuperScrollView.LoopListView2).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);

            instance_of_this_method.InitListView(@itemTotalCount, @onGetItemByIndex, @initParam);

            return __ret;
        }

        static StackObject* NewListViewItem_1(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 2);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            System.String @itemPrefabName = (System.String)typeof(System.String).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 2);
            SuperScrollView.LoopListView2 instance_of_this_method = (SuperScrollView.LoopListView2)typeof(SuperScrollView.LoopListView2).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);

            var result_of_this_method = instance_of_this_method.NewListViewItem(@itemPrefabName);

            object obj_result_of_this_method = result_of_this_method;
            if(obj_result_of_this_method is CrossBindingAdaptorType)
            {    
                return ILIntepreter.PushObject(__ret, __mStack, ((CrossBindingAdaptorType)obj_result_of_this_method).ILInstance);
            }
            return ILIntepreter.PushObject(__ret, __mStack, result_of_this_method);
        }

        static StackObject* SetListItemCount_2(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 3);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            System.Boolean @resetPos = ptr_of_this_method->Value == 1;

            ptr_of_this_method = ILIntepreter.Minus(__esp, 2);
            System.Int32 @itemCount = ptr_of_this_method->Value;

            ptr_of_this_method = ILIntepreter.Minus(__esp, 3);
            SuperScrollView.LoopListView2 instance_of_this_method = (SuperScrollView.LoopListView2)typeof(SuperScrollView.LoopListView2).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);

            instance_of_this_method.SetListItemCount(@itemCount, @resetPos);

            return __ret;
        }

        static StackObject* RefreshAllShownItem_3(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 1);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            SuperScrollView.LoopListView2 instance_of_this_method = (SuperScrollView.LoopListView2)typeof(SuperScrollView.LoopListView2).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);

            instance_of_this_method.RefreshAllShownItem();

            return __ret;
        }

        static StackObject* get_ItemTotalCount_4(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 1);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            SuperScrollView.LoopListView2 instance_of_this_method = (SuperScrollView.LoopListView2)typeof(SuperScrollView.LoopListView2).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);

            var result_of_this_method = instance_of_this_method.ItemTotalCount;

            __ret->ObjectType = ObjectTypes.Integer;
            __ret->Value = result_of_this_method;
            return __ret + 1;
        }

        static StackObject* MovePanelToItemIndex_5(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 3);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            System.Single @offset = *(float*)&ptr_of_this_method->Value;

            ptr_of_this_method = ILIntepreter.Minus(__esp, 2);
            System.Int32 @itemIndex = ptr_of_this_method->Value;

            ptr_of_this_method = ILIntepreter.Minus(__esp, 3);
            SuperScrollView.LoopListView2 instance_of_this_method = (SuperScrollView.LoopListView2)typeof(SuperScrollView.LoopListView2).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);

            instance_of_this_method.MovePanelToItemIndex(@itemIndex, @offset);

            return __ret;
        }

        static StackObject* GetShownItemByItemIndex_6(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 2);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            System.Int32 @itemIndex = ptr_of_this_method->Value;

            ptr_of_this_method = ILIntepreter.Minus(__esp, 2);
            SuperScrollView.LoopListView2 instance_of_this_method = (SuperScrollView.LoopListView2)typeof(SuperScrollView.LoopListView2).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);

            var result_of_this_method = instance_of_this_method.GetShownItemByItemIndex(@itemIndex);

            object obj_result_of_this_method = result_of_this_method;
            if(obj_result_of_this_method is CrossBindingAdaptorType)
            {    
                return ILIntepreter.PushObject(__ret, __mStack, ((CrossBindingAdaptorType)obj_result_of_this_method).ILInstance);
            }
            return ILIntepreter.PushObject(__ret, __mStack, result_of_this_method);
        }

        static StackObject* GetFirstShownItemIndexAndOffset_7(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 1);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            SuperScrollView.LoopListView2 instance_of_this_method = (SuperScrollView.LoopListView2)typeof(SuperScrollView.LoopListView2).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);

            var result_of_this_method = instance_of_this_method.GetFirstShownItemIndexAndOffset();

            return ILIntepreter.PushObject(__ret, __mStack, result_of_this_method);
        }

        static StackObject* get_ArrangeType_8(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 1);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            SuperScrollView.LoopListView2 instance_of_this_method = (SuperScrollView.LoopListView2)typeof(SuperScrollView.LoopListView2).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);

            var result_of_this_method = instance_of_this_method.ArrangeType;

            return ILIntepreter.PushObject(__ret, __mStack, result_of_this_method);
        }



    }
}
