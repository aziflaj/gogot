RSpec.describe 'Reverting time' do
  let(:command_result) { `gogot time-machine #{hash} #{path}` }
  let(:path) { './testdir/file11.txt' }

  around do |example|
    `gogot init kewl-projekt`
    Dir.chdir('kewl-projekt') do
      # Commit 1
      File.open('file1.txt', 'w+') { |f| f.write('Test File 1') }

      Dir.mkdir('testdir')
      File.open('testdir/file11.txt', 'w+') { |f| f.write('Test File 11 in testdir') }
      File.open('testdir/file12.txt', 'w+') { |f| f.write('Test File 12 in testdir') }
      File.open('testdir/file13.txt', 'w+') { |f| f.write('Test File 13 in testdir') }

      Dir.mkdir('testdir/tester-dir')
      File.open('testdir/tester-dir/file11.txt', 'w+') { |f| f.write('Test File 11 in tester-dir') }
      File.open('testdir/tester-dir/file12.txt', 'w+') { |f| f.write('Test File 12 in tester-dir') }
      File.open('testdir/tester-dir/file13.txt', 'w+') { |f| f.write('Test File 13 in tester-dir') }

      Dir.mkdir('testdir2')
      File.open('testdir2/file21.txt', 'w+') { |f| f.write('Test File 21 in testdir2') }
      File.open('testdir2/file22.txt', 'w+') { |f| f.write('Test File 22 in testdir2') }
      File.open('testdir2/file23.txt', 'w+') { |f| f.write('Test File 23 in testdir2') }

      `gogot add .`
      `gogot commit Commit 1`

      # Commit 2
      File.open('file1.txt', 'a') { |f| f.write('Hello') }
      File.open('testdir/file11.txt', 'a') { |f| f.write('Hello') }

      `gogot add .`
      `gogot commit Commit 2`

      example.run
    end
    `rm -rf kewl-projekt`
  end

  context 'without params' do
    let(:path) { '' }

    it 'prints usage information' do
      expect(command_result).to include('Usage: gogot time-machine <commit-id> <file-path>')
    end
  end

  it do
    log_result = `gogot log`
    commits = log_result.split("\n")[1..].map { |line| line.split('    ') }
    result = `gogot time-machine #{commits.first[0]} #{path}`

    expect(result).to include('Test File 11 in testdir')
    expect(result).not_to include('Hello')

    current_content = File.read(path)
    expect(current_content).to include('Hello')
  end
end
