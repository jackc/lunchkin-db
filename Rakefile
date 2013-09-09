begin
  require 'bundler'
  Bundler.setup
rescue LoadError
  puts 'You must `gem install bundler` and `bundle install` to run rake tasks'
end

require 'fileutils'
require 'rspec/core/rake_task'

file 'players_index.go' => 'players_index.gst' do |t|
  t.prerequisites.each do |source|
    sh "gst #{source} | gofmt > #{t.name}"
  end
end

file 'header.go' => 'header.gst' do |t|
  t.prerequisites.each do |source|
    sh "gst #{source} | gofmt > #{t.name}"
  end
end

file 'footer.go' => 'footer.gst' do |t|
  t.prerequisites.each do |source|
    sh "gst #{source} | gofmt > #{t.name}"
  end
end

file 'standings.go' => 'standings.gst' do |t|
  t.prerequisites.each do |source|
    sh "gst #{source} | gofmt > #{t.name}"
  end
end

file 'games_index.go' => 'games_index.gst' do |t|
  t.prerequisites.each do |source|
    sh "gst #{source} | gofmt > #{t.name}"
  end
end

file 'games_new.go' => 'games_new.gst' do |t|
  t.prerequisites.each do |source|
    sh "gst #{source} | gofmt > #{t.name}"
  end
end

file 'lunchkin-db' => ['handlers.go', 'main.go', 'migrations.go', 'players_index.go', 'header.go', 'footer.go', 'standings.go', 'games_index.go', 'games_new.go'] do |t|
  sh 'go build'
end

desc 'Build lunchkin-db'
task build: 'lunchkin-db'

desc 'Run lunchkin-db server'
task server: 'lunchkin-db' do
  exec './lunchkin-db'
end

task spec_server: :build do
  FileUtils.mkdir_p 'tmp/spec/server'
  FileUtils.touch 'tmp/spec/server/stdout.log'
  FileUtils.touch 'tmp/spec/server/stderr.log'
  pid = Process.spawn './lunchkin-db -config=config.test.yml',
    out: 'tmp/spec/server/stdout.log',
    err: 'tmp/spec/server/stderr.log'
  at_exit { Process.kill 'TERM', pid }
end

RSpec::Core::RakeTask.new(:spec)
task spec: :spec_server

task :default => :spec
