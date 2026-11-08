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
    unsafe class TouchEventCom_Binding
    {
        public static void Register(ILRuntime.Runtime.Enviorment.AppDomain app)
        {
            BindingFlags flag = BindingFlags.Public | BindingFlags.Instance | BindingFlags.Static | BindingFlags.DeclaredOnly;
            MethodBase method;
            Type[] args;
            Type type = typeof(global::TouchEventCom);
            args = new Type[]{};
            method = type.GetMethod("RemoveAllListeners", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, RemoveAllListeners_0);
            args = new Type[]{typeof(UnityEngine.Events.UnityAction<UnityEngine.EventSystems.PointerEventData>)};
            method = type.GetMethod("AddDragBeginListener", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, AddDragBeginListener_1);
            args = new Type[]{typeof(UnityEngine.Events.UnityAction<UnityEngine.EventSystems.PointerEventData>)};
            method = type.GetMethod("AddDragEndListener", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, AddDragEndListener_2);
            args = new Type[]{typeof(UnityEngine.Events.UnityAction<UnityEngine.EventSystems.PointerEventData>)};
            method = type.GetMethod("AddDragListener", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, AddDragListener_3);
            args = new Type[]{typeof(UnityEngine.Events.UnityAction<UnityEngine.EventSystems.PointerEventData>)};
            method = type.GetMethod("AddPointerDownListener", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, AddPointerDownListener_4);
            args = new Type[]{typeof(UnityEngine.Events.UnityAction<UnityEngine.EventSystems.PointerEventData>)};
            method = type.GetMethod("AddClickListener", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, AddClickListener_5);
            args = new Type[]{typeof(UnityEngine.Events.UnityAction<UnityEngine.EventSystems.PointerEventData>)};
            method = type.GetMethod("AddPointerUpListener", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, AddPointerUpListener_6);


        }


        static StackObject* RemoveAllListeners_0(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 1);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            global::TouchEventCom instance_of_this_method = (global::TouchEventCom)typeof(global::TouchEventCom).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);

            instance_of_this_method.RemoveAllListeners();

            return __ret;
        }

        static StackObject* AddDragBeginListener_1(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 2);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            UnityEngine.Events.UnityAction<UnityEngine.EventSystems.PointerEventData> @action = (UnityEngine.Events.UnityAction<UnityEngine.EventSystems.PointerEventData>)typeof(UnityEngine.Events.UnityAction<UnityEngine.EventSystems.PointerEventData>).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)8);
            __intp.Free(ptr_of_this_method);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 2);
            global::TouchEventCom instance_of_this_method = (global::TouchEventCom)typeof(global::TouchEventCom).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);

            instance_of_this_method.AddDragBeginListener(@action);

            return __ret;
        }

        static StackObject* AddDragEndListener_2(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 2);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            UnityEngine.Events.UnityAction<UnityEngine.EventSystems.PointerEventData> @action = (UnityEngine.Events.UnityAction<UnityEngine.EventSystems.PointerEventData>)typeof(UnityEngine.Events.UnityAction<UnityEngine.EventSystems.PointerEventData>).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)8);
            __intp.Free(ptr_of_this_method);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 2);
            global::TouchEventCom instance_of_this_method = (global::TouchEventCom)typeof(global::TouchEventCom).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);

            instance_of_this_method.AddDragEndListener(@action);

            return __ret;
        }

        static StackObject* AddDragListener_3(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 2);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            UnityEngine.Events.UnityAction<UnityEngine.EventSystems.PointerEventData> @action = (UnityEngine.Events.UnityAction<UnityEngine.EventSystems.PointerEventData>)typeof(UnityEngine.Events.UnityAction<UnityEngine.EventSystems.PointerEventData>).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)8);
            __intp.Free(ptr_of_this_method);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 2);
            global::TouchEventCom instance_of_this_method = (global::TouchEventCom)typeof(global::TouchEventCom).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);

            instance_of_this_method.AddDragListener(@action);

            return __ret;
        }

        static StackObject* AddPointerDownListener_4(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 2);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            UnityEngine.Events.UnityAction<UnityEngine.EventSystems.PointerEventData> @action = (UnityEngine.Events.UnityAction<UnityEngine.EventSystems.PointerEventData>)typeof(UnityEngine.Events.UnityAction<UnityEngine.EventSystems.PointerEventData>).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)8);
            __intp.Free(ptr_of_this_method);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 2);
            global::TouchEventCom instance_of_this_method = (global::TouchEventCom)typeof(global::TouchEventCom).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);

            instance_of_this_method.AddPointerDownListener(@action);

            return __ret;
        }

        static StackObject* AddClickListener_5(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 2);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            UnityEngine.Events.UnityAction<UnityEngine.EventSystems.PointerEventData> @action = (UnityEngine.Events.UnityAction<UnityEngine.EventSystems.PointerEventData>)typeof(UnityEngine.Events.UnityAction<UnityEngine.EventSystems.PointerEventData>).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)8);
            __intp.Free(ptr_of_this_method);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 2);
            global::TouchEventCom instance_of_this_method = (global::TouchEventCom)typeof(global::TouchEventCom).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);

            instance_of_this_method.AddClickListener(@action);

            return __ret;
        }

        static StackObject* AddPointerUpListener_6(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 2);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            UnityEngine.Events.UnityAction<UnityEngine.EventSystems.PointerEventData> @action = (UnityEngine.Events.UnityAction<UnityEngine.EventSystems.PointerEventData>)typeof(UnityEngine.Events.UnityAction<UnityEngine.EventSystems.PointerEventData>).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)8);
            __intp.Free(ptr_of_this_method);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 2);
            global::TouchEventCom instance_of_this_method = (global::TouchEventCom)typeof(global::TouchEventCom).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)0);
            __intp.Free(ptr_of_this_method);

            instance_of_this_method.AddPointerUpListener(@action);

            return __ret;
        }



    }
}
