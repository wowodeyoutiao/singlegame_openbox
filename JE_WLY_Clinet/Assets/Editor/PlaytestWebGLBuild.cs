using System;
using System.Collections.Generic;
using System.IO;
using UnityEditor;
using UnityEditor.Build.Reporting;
using UnityEngine;
using YooAsset.Editor;

/// <summary>
/// wly 手机网页试玩构建脚本（batchmode 可调用）
///   -executeMethod PlaytestWebGLBuild.BuildBundlesOnly
///   -executeMethod PlaytestWebGLBuild.BuildPlayerOnly
///   -executeMethod PlaytestWebGLBuild.BuildAll
/// </summary>
public static class PlaytestWebGLBuild
{
    const string WebOutputDir = @"D:\wly_run\web";
    const string LogPath = @"D:\wly_run\webgl_build.log";

    static readonly string[] Packages = { "Main", "AddOn1" };

    static void Log(string msg)
    {
        string line = "[WLY-WEBGL] " + msg;
        Debug.Log(line);
        try
        {
            Directory.CreateDirectory(Path.GetDirectoryName(LogPath));
            File.AppendAllText(LogPath, line + "\n");
        }
        catch { }
    }

    // ------------------------------------------------------------------ 平台设置
    static void ConfigurePlayerSettings()
    {
        // 注意：这里只做检查与日志，不再在运行时修改 PlayerSettings，
        // 因为改动宏/后端会触发重新编译和域重载，导致 executeMethod 被中断。

        // 修复：PlayerSettings 里 WebGL 模板被设成 "Mobile"，但 Unity 2022.3.62 的
        // WebGLSupport\BuildTools\WebGLTemplates 下只有 Default/Minimal/PWA/WebGLIncludes，
        // 会报 "Invalid WebGL template path ... WebGLTemplates\Mobile" 导致构建失败。
        // 强制改用内置 Default 模板。
        if (PlayerSettings.WebGL.template != "APPLICATION:Default")
        {
            Log("template 修正: [" + PlayerSettings.WebGL.template + "] -> APPLICATION:Default");
            PlayerSettings.WebGL.template = "APPLICATION:Default";
            AssetDatabase.SaveAssets();
        }

        Log("template=" + PlayerSettings.WebGL.template
            + " compression=" + PlayerSettings.WebGL.compressionFormat
            + " mem=" + PlayerSettings.WebGL.memorySize
            + " backend=" + PlayerSettings.GetScriptingBackend(BuildTargetGroup.WebGL)
            + " defines=" + PlayerSettings.GetScriptingDefineSymbolsForGroup(BuildTargetGroup.WebGL));
    }

    // ------------------------------------------------------------------ 资源包
    static bool BuildOnePackage(string packageName)
    {
        // 用内建管线（BBP）而不是 SBP：SBP 在打 WebGL 包时会去编译 player 脚本，
        // 本机 IL2CPP 环境下必然失败（ErrorCode300），内建管线不需要这一步。
        var parameters = new BuiltinBuildParameters
        {
            BuildTarget = BuildTarget.WebGL,
            BuildPipeline = EBuildPipeline.BuiltinBuildPipeline.ToString(),
            BuildBundleType = (int)YooAsset.EBuildBundleType.AssetBundle,
            BuildOutputRoot = Path.Combine(Directory.GetCurrentDirectory(), "Bundles"),
            BuildinFileRoot = Path.Combine(Directory.GetCurrentDirectory(), "Assets", "StreamingAssets", "yoo"),
            PackageName = packageName,
            PackageVersion = DateTime.Now.ToString("yyyy-MM-dd-HHmm"),
            PackageNote = "playtest-webgl",
            CompressOption = ECompressOption.LZ4,
            FileNameStyle = YooAsset.EFileNameStyle.HashName,
            BuildinFileCopyOption = EBuildinFileCopyOption.ClearAndCopyAll,
            VerifyBuildingResult = false,
            EnableSharePackRule = true,
            ClearBuildCacheFiles = false,
            UseAssetDependencyDB = true,
        };

        Log($"开始构建资源包 {packageName} -> {parameters.BuildOutputRoot}");
        var pipeline = new BuiltinBuildPipeline();
        YooAsset.Editor.BuildResult result = pipeline.Run(parameters, true);
        if (result.Success)
        {
            Log($"资源包 {packageName} 构建成功: {result.OutputPackageDirectory}");
            return true;
        }
        Log($"资源包 {packageName} 构建失败: {result.FailedTask} / {result.ErrorInfo}");
        return false;
    }

    public static void BuildBundles()
    {
        ConfigurePlayerSettings();
        bool ok = true;
        foreach (var p in Packages)
        {
            if (!BuildOnePackage(p)) ok = false;
        }
        Log("BuildBundles done ok=" + ok);
        if (!ok) throw new Exception("bundle build failed");
    }

    // ------------------------------------------------------------------ WebGL 包体
    static string[] GetScenes()
    {
        var list = new List<string>();
        foreach (var s in EditorBuildSettings.scenes)
            if (s.enabled) list.Add(s.path);
        if (list.Count == 0) list.Add("Assets/Init.unity");
        return list.ToArray();
    }

    public static void BuildPlayer()
    {
        ConfigurePlayerSettings();

        // 试玩用：关闭 Brotli 压缩，避免本地 python 网关需要返回 Content-Encoding: br
        PlayerSettings.WebGL.compressionFormat = WebGLCompressionFormat.Disabled;
        Log("compression set to Disabled");

        if (Directory.Exists(WebOutputDir))
        {
            try { Directory.Delete(WebOutputDir, true); } catch (Exception e) { Log("清目录失败 " + e.Message); }
        }
        Directory.CreateDirectory(WebOutputDir);

        var scenes = GetScenes();
        Log("scenes = " + string.Join(", ", scenes));

        var options = new BuildPlayerOptions
        {
            scenes = scenes,
            locationPathName = WebOutputDir,
            target = BuildTarget.WebGL,
            targetGroup = BuildTargetGroup.WebGL,
            options = BuildOptions.None,
        };

        var report = BuildPipeline.BuildPlayer(options);
        var summary = report.summary;
        Log($"WEBGL_BUILD result={summary.result} errors={summary.totalErrors} size={summary.totalSize} time={summary.totalTime}");
        if (summary.result != UnityEditor.Build.Reporting.BuildResult.Succeeded)
            throw new Exception("WebGL build failed: " + summary.result);
    }

    // ------------------------------------------------------------------ 入口
    public static void BuildBundlesOnly()
    {
        BuildBundles();
        EditorApplication.Exit(0);
    }

    public static void BuildPlayerOnly()
    {
        BuildPlayer();
        Log("WEBGL_DONE");
        EditorApplication.Exit(0);
    }

    public static void BuildAll()
    {
        BuildBundles();
        BuildPlayer();
        Log("WEBGL_DONE");
        EditorApplication.Exit(0);
    }
}
