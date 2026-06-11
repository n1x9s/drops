require "xcodeproj"

root = File.expand_path("..", __dir__)
ios_root = File.join(root, "ios", "SecondBrain")
project_path = File.join(ios_root, "SecondBrain.xcodeproj")

project = Xcodeproj::Project.new(project_path)
project.root_object.attributes["LastSwiftUpdateCheck"] = "2650"
project.root_object.attributes["LastUpgradeCheck"] = "2650"

app_target = project.new_target(:application, "SecondBrain", :ios, "26.0")
app_target.product_name = "SecondBrain"
app_target.product_reference.name = "SecondBrain.app"

sources_phase = app_target.source_build_phase
resources_phase = app_target.resources_build_phase

main_group = project.main_group
main_group.set_source_tree("<group>")
main_group.path = nil

app_group = main_group.new_group("SecondBrain", ios_root)

swift_files = Dir[File.join(ios_root, "**", "*.swift")]
  .reject { |path| path.include?("/Widgets/") }
  .sort

swift_files.each do |path|
  ref = app_group.new_file(path)
  sources_phase.add_file_reference(ref)
end

info_ref = app_group.new_file(File.join(ios_root, "Info.plist"))

app_target.build_configurations.each do |config|
  config.build_settings["ASSETCATALOG_COMPILER_APPICON_NAME"] = "AppIcon"
  config.build_settings["CODE_SIGN_STYLE"] = "Automatic"
  config.build_settings["CURRENT_PROJECT_VERSION"] = "1"
  config.build_settings["DEVELOPMENT_TEAM"] = ""
  config.build_settings["GENERATE_INFOPLIST_FILE"] = "NO"
  config.build_settings["INFOPLIST_FILE"] = "Info.plist"
  config.build_settings["IPHONEOS_DEPLOYMENT_TARGET"] = "26.0"
  config.build_settings["MARKETING_VERSION"] = "1.0"
  config.build_settings["PRODUCT_BUNDLE_IDENTIFIER"] = "com.secondbrain.app"
  config.build_settings["PRODUCT_NAME"] = "$(TARGET_NAME)"
  config.build_settings["SDKROOT"] = "iphoneos"
  config.build_settings["SUPPORTED_PLATFORMS"] = "iphoneos iphonesimulator"
  config.build_settings["SUPPORTS_MACCATALYST"] = "NO"
  config.build_settings["SWIFT_EMIT_LOC_STRINGS"] = "YES"
  config.build_settings["SWIFT_VERSION"] = "6.0"
  config.build_settings["TARGETED_DEVICE_FAMILY"] = "1,2"
end

project.build_configurations.each do |config|
  config.build_settings["ENABLE_USER_SCRIPT_SANDBOXING"] = "YES"
  config.build_settings["SWIFT_VERSION"] = "6.0"
end

project.save
