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
    unsafe class Cysharp_Threading_Tasks_UniTask_Binding
    {
        public static void Register(ILRuntime.Runtime.Enviorment.AppDomain app)
        {
            BindingFlags flag = BindingFlags.Public | BindingFlags.Instance | BindingFlags.Static | BindingFlags.DeclaredOnly;
            MethodBase method;
            Type[] args;
            Type type = typeof(Cysharp.Threading.Tasks.UniTask);
            args = new Type[]{};
            method = type.GetMethod("GetAwaiter", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, GetAwaiter_0);
            args = new Type[]{typeof(System.Func<System.Boolean>), typeof(Cysharp.Threading.Tasks.PlayerLoopTiming), typeof(System.Threading.CancellationToken)};
            method = type.GetMethod("WaitUntil", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, WaitUntil_1);
            args = new Type[]{typeof(System.Int32), typeof(System.Boolean), typeof(Cysharp.Threading.Tasks.PlayerLoopTiming), typeof(System.Threading.CancellationToken)};
            method = type.GetMethod("Delay", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, Delay_2);
            args = new Type[]{typeof(System.TimeSpan), typeof(Cysharp.Threading.Tasks.DelayType), typeof(Cysharp.Threading.Tasks.PlayerLoopTiming), typeof(System.Threading.CancellationToken)};
            method = type.GetMethod("Delay", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, Delay_3);
            args = new Type[]{typeof(Cysharp.Threading.Tasks.PlayerLoopTiming)};
            method = type.GetMethod("Yield", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, Yield_4);
            args = new Type[]{typeof(System.Threading.CancellationToken)};
            method = type.GetMethod("SwitchToMainThread", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, SwitchToMainThread_5);
            args = new Type[]{};
            method = type.GetMethod("Yield", flag, null, args, null);
            app.RegisterCLRMethodRedirection(method, Yield_6);

            app.RegisterCLRCreateDefaultInstance(type, () => new Cysharp.Threading.Tasks.UniTask());


        }

        static void WriteBackInstance(ILRuntime.Runtime.Enviorment.AppDomain __domain, StackObject* ptr_of_this_method, AutoList __mStack, ref Cysharp.Threading.Tasks.UniTask instance_of_this_method)
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
                        var instance_of_arrayReference = __mStack[ptr_of_this_method->Value] as Cysharp.Threading.Tasks.UniTask[];
                        instance_of_arrayReference[ptr_of_this_method->ValueLow] = instance_of_this_method;
                    }
                    break;
            }
        }

        static StackObject* GetAwaiter_0(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 1);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            ptr_of_this_method = ILIntepreter.GetObjectAndResolveReference(ptr_of_this_method);
            Cysharp.Threading.Tasks.UniTask instance_of_this_method = (Cysharp.Threading.Tasks.UniTask)typeof(Cysharp.Threading.Tasks.UniTask).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)16);

            var result_of_this_method = instance_of_this_method.GetAwaiter();

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            WriteBackInstance(__domain, ptr_of_this_method, __mStack, ref instance_of_this_method);

            __intp.Free(ptr_of_this_method);
            return ILIntepreter.PushObject(__ret, __mStack, result_of_this_method);
        }

        static StackObject* WaitUntil_1(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 3);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            System.Threading.CancellationToken @cancellationToken = (System.Threading.CancellationToken)typeof(System.Threading.CancellationToken).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)16);
            __intp.Free(ptr_of_this_method);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 2);
            Cysharp.Threading.Tasks.PlayerLoopTiming @timing = (Cysharp.Threading.Tasks.PlayerLoopTiming)typeof(Cysharp.Threading.Tasks.PlayerLoopTiming).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)20);
            __intp.Free(ptr_of_this_method);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 3);
            System.Func<System.Boolean> @predicate = (System.Func<System.Boolean>)typeof(System.Func<System.Boolean>).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)8);
            __intp.Free(ptr_of_this_method);


            var result_of_this_method = Cysharp.Threading.Tasks.UniTask.WaitUntil(@predicate, @timing, @cancellationToken);

            return ILIntepreter.PushObject(__ret, __mStack, result_of_this_method);
        }

        static StackObject* Delay_2(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 4);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            System.Threading.CancellationToken @cancellationToken = (System.Threading.CancellationToken)typeof(System.Threading.CancellationToken).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)16);
            __intp.Free(ptr_of_this_method);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 2);
            Cysharp.Threading.Tasks.PlayerLoopTiming @delayTiming = (Cysharp.Threading.Tasks.PlayerLoopTiming)typeof(Cysharp.Threading.Tasks.PlayerLoopTiming).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)20);
            __intp.Free(ptr_of_this_method);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 3);
            System.Boolean @ignoreTimeScale = ptr_of_this_method->Value == 1;

            ptr_of_this_method = ILIntepreter.Minus(__esp, 4);
            System.Int32 @millisecondsDelay = ptr_of_this_method->Value;


            var result_of_this_method = Cysharp.Threading.Tasks.UniTask.Delay(@millisecondsDelay, @ignoreTimeScale, @delayTiming, @cancellationToken);

            return ILIntepreter.PushObject(__ret, __mStack, result_of_this_method);
        }

        static StackObject* Delay_3(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 4);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            System.Threading.CancellationToken @cancellationToken = (System.Threading.CancellationToken)typeof(System.Threading.CancellationToken).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)16);
            __intp.Free(ptr_of_this_method);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 2);
            Cysharp.Threading.Tasks.PlayerLoopTiming @delayTiming = (Cysharp.Threading.Tasks.PlayerLoopTiming)typeof(Cysharp.Threading.Tasks.PlayerLoopTiming).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)20);
            __intp.Free(ptr_of_this_method);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 3);
            Cysharp.Threading.Tasks.DelayType @delayType = (Cysharp.Threading.Tasks.DelayType)typeof(Cysharp.Threading.Tasks.DelayType).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)20);
            __intp.Free(ptr_of_this_method);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 4);
            System.TimeSpan @delayTimeSpan = (System.TimeSpan)typeof(System.TimeSpan).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)16);
            __intp.Free(ptr_of_this_method);


            var result_of_this_method = Cysharp.Threading.Tasks.UniTask.Delay(@delayTimeSpan, @delayType, @delayTiming, @cancellationToken);

            return ILIntepreter.PushObject(__ret, __mStack, result_of_this_method);
        }

        static StackObject* Yield_4(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 1);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            Cysharp.Threading.Tasks.PlayerLoopTiming @timing = (Cysharp.Threading.Tasks.PlayerLoopTiming)typeof(Cysharp.Threading.Tasks.PlayerLoopTiming).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)20);
            __intp.Free(ptr_of_this_method);


            var result_of_this_method = Cysharp.Threading.Tasks.UniTask.Yield(@timing);

            return ILIntepreter.PushObject(__ret, __mStack, result_of_this_method);
        }

        static StackObject* SwitchToMainThread_5(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* ptr_of_this_method;
            StackObject* __ret = ILIntepreter.Minus(__esp, 1);

            ptr_of_this_method = ILIntepreter.Minus(__esp, 1);
            System.Threading.CancellationToken @cancellationToken = (System.Threading.CancellationToken)typeof(System.Threading.CancellationToken).CheckCLRTypes(StackObject.ToObject(ptr_of_this_method, __domain, __mStack), (CLR.Utils.Extensions.TypeFlags)16);
            __intp.Free(ptr_of_this_method);


            var result_of_this_method = Cysharp.Threading.Tasks.UniTask.SwitchToMainThread(@cancellationToken);

            return ILIntepreter.PushObject(__ret, __mStack, result_of_this_method);
        }

        static StackObject* Yield_6(ILIntepreter __intp, StackObject* __esp, AutoList __mStack, CLRMethod __method, bool isNewObj)
        {
            ILRuntime.Runtime.Enviorment.AppDomain __domain = __intp.AppDomain;
            StackObject* __ret = ILIntepreter.Minus(__esp, 0);


            var result_of_this_method = Cysharp.Threading.Tasks.UniTask.Yield();

            return ILIntepreter.PushObject(__ret, __mStack, result_of_this_method);
        }



    }
}
