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
    unsafe class SuperScrollView_LoopGridView_Binding
    {
        public static void Register(ILRuntime.Runtime.Enviorment.AppDomain app)
        {
            BindingFlags flag = BindingFlags.Public | BindingFlags.Instance | BindingFlags.Static | BindingFlags.DeclaredOnly;
            MethodBase method;
            Type[] args;
            Type type = typeof(SuperScrollView.LoopGridView);
            args = new Type[]{typeof(System.String)};
            method = type.GetMethod("NewListViewItem", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, NewListViewItem_0);
            args = new Type[]{typeof(System.Int32), typeof(System.Boolean)};
            method = type.GetMethod("SetListItemCount", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, SetListItemCount_1);
            args = new Type[]{};
            method = type.GetMethod("RefreshAllShownItem", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, RefreshAllShownItem_2);
            args = new Type[]{typeof(System.Int32), typeof(System.Single), typeof(System.Single)};
            method = type.GetMethod("MovePanelToItemByIndex", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, MovePanelToItemByIndex_3);
            args = new Type[]{typeof(System.Int32), typeof(System.Func<SuperScrollView.LoopGridView, System.Int32, System.Int32, System.Int32, SuperScrollView.LoopGridViewItem>), typeof(SuperScrollView.LoopGridViewSettingParam), typeof(SuperScrollView.LoopGridViewInitParam)};
            method = type.GetMethod("InitGridView", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, InitGridView_4);
            args = new Type[]{};
            method = type.GetMethod("get_ItemTotalCount", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, get_ItemTotalCount_5);
            args = new Type[]{typeof(System.Int32), typeof(System.Int32), typeof(System.Single), typeof(System.Single)};
            method = type.GetMethod("MovePanelToItemByRowColumn", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, MovePanelToItemByRowColumn_6);


        }


        static StackObject* NewListViewItem_0(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 2);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            System.String @itemPrefabName = (System.String)typeof(System.String).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 2);
            SuperScrollView.LoopGridView instance_of_this_method = (SuperScrollView.LoopGridView)typeof(SuperScrollView.LoopGridView).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);

            var result_of_this_method = instance_of_this_method.NewListViewItem(@itemPrefabName);

            object obj_result_of_this_method = result_of_this_method;
            if(obj_result_of_this_method is CrossBindingAdaptorType)
            {    
                return ILIntepreter.PushObject(__ret, __mStack, ((CrossBindingAdaptorType)obj_result_of_this_method).ILInstance);
            }
            return ILIntepreter.PushObject(__ret, __mStack, result_of_this_method);
        }

        static StackObject* SetListItemCount_1(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 3);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            System.Boolean @resetPos = ptr_of_this_method->Value == 1;

            ptr_of_this_method = ILIntepreter.Minus(__esp, 2);
            System.Int32 @itemCount = ptr_of_this_method->Value;

            ptr_of_this_method = ILIntepreter.Minus(__esp, 3);
            SuperScrollView.LoopGridView instance_of_this_method = (SuperScrollView.LoopGridView)typeof(SuperScrollView.LoopGridView).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);

            instance_of_this_method.SetListItemCount(@itemCount, @resetPos);

            return __ret;
        }

        static StackObject* RefreshAllShownItem_2(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 1);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            SuperScrollView.LoopGridView instance_of_this_method = (SuperScrollView.LoopGridView)typeof(SuperScrollView.LoopGridView).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);

            instance_of_this_method.RefreshAllShownItem();

            return __ret;
        }

        static StackObject* MovePanelToItemByIndex_3(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 4);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            System.Single @offsetY = *(float*)&ptr_of_this_method->Value;

            ptr_of_this_method = ILIntepreter.Minus(__esp, 2);
            System.Single @offsetX = *(float*)&ptr_of_this_method->Value;

            ptr_of_this_method = ILIntepreter.Minus(__esp, 3);
            System.Int32 @itemIndex = ptr_of_this_method->Value;

            ptr_of_this_method = ILIntepreter.Minus(__esp, 4);
            SuperScrollView.LoopGridView instance_of_this_method = (SuperScrollView.LoopGridView)typeof(SuperScrollView.LoopGridView).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);

            instance_of_this_method.MovePanelToItemByIndex(@itemIndex, @offsetX, @offsetY);

            return __ret;
        }

        static StackObject* InitGridView_4(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 5);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            SuperScrollView.LoopGridViewInitParam @initParam = (SuperScrollView.LoopGridViewInitParam)typeof(SuperScrollView.LoopGridViewInitParam).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 2);
            SuperScrollView.LoopGridViewSettingParam @settingParam = (SuperScrollView.LoopGridViewSettingParam)typeof(SuperScrollView.LoopGridViewSettingParam).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 3);
            System.Func<SuperScrollView.LoopGridView, System.Int32, System.Int32, System.Int32, SuperScrollView.LoopGridViewItem> @onGetItemByRowColumn = (System.Func<SuperScrollView.LoopGridView, System.Int32, System.Int32, System.Int32, SuperScrollView.LoopGridViewItem>)typeof(System.Func<SuperScrollView.LoopGridView, System.Int32, System.Int32, System.Int32, SuperScrollView.LoopGridViewItem>).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)8);
            __intp.Free(ptr_of_this_method);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 4);
            System.Int32 @itemTotalCount = ptr_of_this_method->Value;

            ptr_of_this_method = ILIntepreter.Minus(__esp, 5);
            SuperScrollView.LoopGridView instance_of_this_method = (SuperScrollView.LoopGridView)typeof(SuperScrollView.LoopGridView).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);

            instance_of_this_method.InitGridView(@itemTotalCount, @onGetItemByRowColumn, @settingParam, @initParam);

            return __ret;
        }

        static StackObject* get_ItemTotalCount_5(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 1);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            SuperScrollView.LoopGridView instance_of_this_method = (SuperScrollView.LoopGridView)typeof(SuperScrollView.LoopGridView).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);

            var result_of_this_method = instance_of_this_method.ItemTotalCount;

            __ret->ObjectType = ObjectTypes.Integer;
            __ret->Value = result_of_this_method;
            return __ret + 1;
        }

        static StackObject* MovePanelToItemByRowColumn_6(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 5);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            System.Single @offsetY = *(float*)&ptr_of_this_method->Value;

            ptr_of_this_method = ILIntepreter.Minus(__esp, 2);
            System.Single @offsetX = *(float*)&ptr_of_this_method->Value;

            ptr_of_this_method = ILIntepreter.Minus(__esp, 3);
            System.Int32 @column = ptr_of_this_method->Value;

            ptr_of_this_method = ILIntepreter.Minus(__esp, 4);
            System.Int32 @row = ptr_of_this_method->Value;

            ptr_of_this_method = ILIntepreter.Minus(__esp, 5);
            SuperScrollView.LoopGridView instance_of_this_method = (SuperScrollView.LoopGridView)typeof(SuperScrollView.LoopGridView).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);

            instance_of_this_method.MovePanelToItemByRowColumn(@row, @column, @offsetX, @offsetY);

            return __ret;
        }



    }
}
